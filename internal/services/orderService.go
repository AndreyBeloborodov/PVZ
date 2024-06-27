package services

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"gitlab.ozon.dev/go/classroom-13/students/Homework/Homework-1/internal/models"
	"gitlab.ozon.dev/go/classroom-13/students/Homework/Homework-1/internal/storage"
	"gitlab.ozon.dev/go/classroom-13/students/Homework/Homework-1/internal/storage/transactor"
	"math"
	"time"
)

type OrderService struct {
	Storage   *storage.OrderStorage
	TxManager *transactor.TransactionManager
}

type ParamsGetOrders struct {
	Limit     int
	OnlyExist bool
}

type ParamsAddOrder struct {
	OrderId     int
	UserId      int
	TimeEnd     time.Time
	Weight      int
	Price       int
	PackageName string
}

func NewOrderService(storage *storage.OrderStorage, txManager *transactor.TransactionManager) *OrderService {
	return &OrderService{
		Storage:   storage,
		TxManager: txManager,
	}
}

func (s *OrderService) AddOrder(ctx context.Context, params ParamsAddOrder) error {

	err := s.TxManager.RunReadCommitted(ctx, func(ctxTX context.Context) error {
		exist, err := s.Storage.IsExistOrder(ctxTX, models.OrderId(params.OrderId))

		if err != nil {
			return errors.New("не удалось принять заказ")
		}

		if exist {
			return errors.New("заказ уже принят")
		}

		if params.TimeEnd.Before(time.Now()) {
			return errors.New("срок хранения в прошлом")
		}

		Package, err := NewPackageService((*storage.PackageStorage)(s.Storage), s.TxManager).GetPackageByName(ctxTX, params.PackageName)

		if err != nil {
			return errors.Wrap(err, "не удалось принять заказ")
		}

		if int(Package.MaxWeight) != -1 && params.Weight > int(Package.MaxWeight) {
			return errors.New("заказ слишком тяжёлый для этой упаковки")
		}

		err = s.Storage.AddOrder(
			ctxTX,
			models.Order{
				OrderId:     models.OrderId(params.OrderId),
				UserId:      models.UserId(params.UserId),
				TimeEnd:     models.TimeEnd(params.TimeEnd),
				Weight:      models.OrderWeight(params.Weight),
				FirstPrice:  models.FirstPrice(params.Price),
				PackageName: models.PackageName(params.PackageName),
				ResultPrice: models.ResultPrice(params.Price + int(Package.Price)),
			})

		if err != nil {
			return errors.New("не удалось принять заказ")
		}

		return nil
	})

	return err
}

func (s *OrderService) ReturnOrderToCourier(ctx context.Context, orderId int) error {

	order, err := s.Storage.GetOrderById(ctx, models.OrderId(orderId))
	if err != nil {
		return errors.Wrap(err, "не удалось вернуть заказ курьеру")
	}

	if time.Now().Before(time.Time(order.TimeEnd)) {
		return errors.New("у заказа ещё не вышел срок хранения")
	}

	if order.IsGiven {
		return errors.New("заказ уже выдан клиенту")
	}

	if _, err = s.Storage.DeleteOrderById(ctx, models.OrderId(orderId)); err != nil {
		return errors.New("не удалось вернуть заказ курьеру")
	}

	return nil
}

func (s *OrderService) GiveOrderToUser(ctx context.Context, ordersId []int) error {

	err := s.TxManager.RunReadCommitted(ctx, func(ctxTX context.Context) error {

		orders := make([]models.Order, 0)
		for _, orderId := range ordersId {
			order, err := s.Storage.GetOrderById(ctxTX, models.OrderId(orderId))
			if err != nil {
				continue
			}
			orders = append(orders, order)
		}

		for i := 1; i < len(orders); i++ {
			if orders[i].UserId != orders[0].UserId {
				return errors.New("заказы принадлежат не одному клиенту")
			}
		}

		for _, order := range orders {
			if !time.Now().Before(time.Time(order.TimeEnd)) {
				fmt.Println("у заказа с id", int(order.OrderId), "кончился срок хранения")
				continue
			}

			err := s.Storage.GiveOrderById(ctxTX, order.OrderId)
			if err != nil {
				return err
			}
		}

		return nil
	})

	return err
}

func (s *OrderService) GetOrders(ctx context.Context, userId int, params ParamsGetOrders) ([]models.Order, error) {

	orders, err := s.Storage.GetAllOrders(ctx)
	if err != nil {
		return nil, errors.New("не удалось вернуть заказы")
	}

	if params.Limit < 0 {
		params.Limit = math.MaxInt32
	}

	resultOrders := make([]models.Order, 0)
	for i := len(orders) - 1; i >= 0 && len(resultOrders) < params.Limit; i-- {
		order := orders[i]
		if order.UserId == models.UserId(userId) && (!params.OnlyExist || bool(!order.IsGiven)) {
			resultOrders = append(resultOrders, order)
		}
	}

	return resultOrders, nil
}

func (s *OrderService) AcceptReturnOrder(ctx context.Context, orderId, userId int) error {

	err := s.TxManager.RunReadCommitted(ctx, func(ctxTX context.Context) error {

		order, err := s.Storage.GetOrderById(ctxTX, models.OrderId(orderId))
		if err != nil {
			return errors.New("не удалось принять возврат")
		}

		if userId != int(order.UserId) {
			return errors.New("заказ не принадлежит этому клиенту")
		}

		if !order.IsGiven {
			return errors.New("заказ не выдан")
		}

		if order.IsReturned {
			return errors.New("заказ уже возвращён")
		}

		if !time.Now().Before(time.Time(order.TimeGiven).Add(time.Hour * 24 * 2)) {
			return errors.New("время на возврат заказа истекло")
		}

		err = s.Storage.ReturnOrderByOrderAndUserId(ctxTX, models.OrderId(orderId), models.UserId(userId))

		if err != nil {
			return errors.New("не удалось принять возврат")
		}

		return nil
	})

	return err
}

func (s *OrderService) GetReturns(ctx context.Context) ([]models.Order, error) {

	orders, err := s.Storage.GetAllOrders(ctx)
	if err != nil {
		return nil, errors.New("не удалось найти возвраты")
	}

	resultOrders := make([]models.Order, 0)
	for _, order := range orders {
		if order.IsReturned {
			resultOrders = append(resultOrders, order)
		}
	}

	return resultOrders, nil
}

func (s *OrderService) GetReturnsPagination(ctx context.Context) func(count int) []models.Order {

	start, end := 0, 0 // [start, end)

	returns, _ := s.GetReturns(ctx)

	return func(count int) []models.Order {
		start = end
		end += count

		if end > len(returns) {
			end = len(returns)
		}

		return returns[start:end]
	}
}
