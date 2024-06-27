package storage

import (
	"context"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/pkg/errors"
	"gitlab.ozon.dev/go/classroom-13/students/Homework/Homework-1/internal/models"
	"gitlab.ozon.dev/go/classroom-13/students/Homework/Homework-1/internal/storage/transactor"
	"time"
)

type OrderStorage struct {
	Provider transactor.QueryEngineProvider
}

func NewOrderStorage(provider transactor.QueryEngineProvider) *OrderStorage {
	return &OrderStorage{provider}
}

func (s *OrderStorage) AddOrder(ctx context.Context, order models.Order) error {
	db := s.Provider.GetQueryEngine(ctx)

	orderRecord := transform(order)

	query := `INSERT INTO orders VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`

	rows, err := db.Query(ctx, query,
		orderRecord.OrderId,
		orderRecord.UserId,
		orderRecord.TimeEnd,
		orderRecord.IsGiven,
		orderRecord.TimeGiven,
		orderRecord.IsReturned,
		orderRecord.Price,
		orderRecord.Weight,
		orderRecord.PackageName,
		orderRecord.ResultPrice,
	)
	defer rows.Close()

	if err != nil {
		return errors.Wrap(err, "не удалось выполнить запрос")
	}
	return nil
}

func (s *OrderStorage) DeleteOrderById(ctx context.Context, orderId models.OrderId) (models.Order, error) {
	db := s.Provider.GetQueryEngine(ctx)

	// Сначала выбираем запись, чтобы вернуть удалённый заказ
	querySelect := `SELECT * FROM orders WHERE id=$1`

	rows, err := db.Query(ctx, querySelect, int(orderId))
	defer rows.Close()

	if err != nil {
		return models.Order{}, errors.Wrap(err, "не удалось выполнить запрос на получение данных заказа")
	}

	var orderRecord OrderRecord
	if err := pgxscan.ScanOne(&orderRecord, rows); err != nil {
		return models.Order{}, err
	}

	if rows.Err() != nil {
		return models.Order{}, errors.Wrap(rows.Err(), "ошибка при чтении результатов запроса")
	}

	// Затем удаляем запись
	queryDelete := `DELETE FROM orders WHERE id=$1`
	_, err = db.Query(ctx, queryDelete, int(orderId))
	if err != nil {
		return models.Order{}, errors.Wrap(err, "не удалось выполнить запрос на удаление")
	}

	return orderRecord.toDomain(), nil
}

func (s *OrderStorage) GetOrderById(ctx context.Context, orderId models.OrderId) (models.Order, error) {
	db := s.Provider.GetQueryEngine(ctx)

	// Запрос для получения заказа по ID
	query := `SELECT * FROM orders WHERE id=$1`
	rows, err := db.Query(ctx, query, int(orderId))
	defer rows.Close()

	if err != nil {
		return models.Order{}, errors.Wrap(err, "не удалось выполнить запрос")
	}

	var orderRecord OrderRecord
	if err := pgxscan.ScanOne(&orderRecord, rows); err != nil {
		return models.Order{}, err
	}

	if rows.Err() != nil {
		return models.Order{}, errors.Wrap(rows.Err(), "ошибка при чтении результатов запроса")
	}

	return orderRecord.toDomain(), nil
}

func (s *OrderStorage) IsExistOrder(ctx context.Context, orderId models.OrderId) (bool, error) {
	db := s.Provider.GetQueryEngine(ctx)

	query := `SELECT EXISTS(SELECT 1 FROM orders WHERE id=$1)`

	rows, err := db.Query(ctx, query, int(orderId))
	defer rows.Close()

	if err != nil {
		return false, errors.Wrap(err, "не удалось выполнить запрос")
	}

	var exists bool
	if rows.Next() {
		if err := rows.Scan(&exists); err != nil {
			return false, errors.Wrap(err, "не удалось получить результат")
		}
	}

	if rows.Err() != nil {
		return false, errors.Wrap(rows.Err(), "ошибка при чтении результатов запроса")
	}

	return exists, nil
}

func (s *OrderStorage) GiveOrderById(ctx context.Context, orderId models.OrderId) error {
	db := s.Provider.GetQueryEngine(ctx)

	// Запрос для обновления флага is_given и установки времени выдачи
	updateQuery := `UPDATE orders SET is_given = true, time_given = $2 WHERE id = $1`

	rows, err := db.Query(ctx, updateQuery, int(orderId), time.Now())
	defer rows.Close()

	if err != nil {
		return errors.Wrap(err, "не удалось выполнить запрос обновления")
	}

	return nil
}

func (s *OrderStorage) ReturnOrderByOrderAndUserId(ctx context.Context, orderId models.OrderId, userId models.UserId) error {
	db := s.Provider.GetQueryEngine(ctx)

	// Запрос для обновления флага is_returned
	updateQuery := `UPDATE orders SET is_returned = true WHERE id = $1 AND user_id = $2`

	rows, err := db.Query(ctx, updateQuery, int(orderId), int(userId))
	defer rows.Close()

	if err != nil {
		return errors.Wrap(err, "не удалось выполнить запрос обновления")
	}

	return nil
}

func (s *OrderStorage) GetAllOrders(ctx context.Context) ([]models.Order, error) {
	db := s.Provider.GetQueryEngine(ctx)

	// Запрос для выборки всех заказов
	query := `SELECT * FROM orders`

	rows, err := db.Query(ctx, query)
	defer rows.Close()

	if err != nil {
		return nil, errors.Wrap(err, "не удалось выполнить запрос получения всех заказов")
	}

	var orderRecords []OrderRecord
	if err := pgxscan.ScanAll(&orderRecords, rows); err != nil {
		return []models.Order{}, err
	}

	if rows.Err() != nil {
		return nil, errors.Wrap(rows.Err(), "ошибка при чтении результатов запроса")
	}

	orders := toDomains(orderRecords)

	return orders, nil
}
