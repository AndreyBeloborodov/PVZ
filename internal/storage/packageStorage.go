package storage

import (
	"context"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/pkg/errors"
	"gitlab.ozon.dev/go/classroom-13/students/Homework/Homework-1/internal/models"
	"gitlab.ozon.dev/go/classroom-13/students/Homework/Homework-1/internal/storage/transactor"
)

type PackageStorage struct {
	Provider transactor.QueryEngineProvider
}

func NewPackageStorage(provider transactor.QueryEngineProvider) *PackageStorage {
	return &PackageStorage{provider}
}

func (s *PackageStorage) GetPackageByName(ctx context.Context, packageName models.PackageName) (models.Package, error) {
	db := s.Provider.GetQueryEngine(ctx)

	// Запрос для получения заказа по ID
	query := `SELECT * FROM packages WHERE name=$1`
	rows, err := db.Query(ctx, query, string(packageName))
	if err != nil {
		return models.Package{}, errors.Wrap(err, "не удалось выполнить запрос")
	}
	defer rows.Close()

	var packageRecord PackageRecord
	if err := pgxscan.ScanOne(&packageRecord, rows); err != nil {
		return models.Package{}, err
	}

	if rows.Err() != nil {
		return models.Package{}, errors.Wrap(rows.Err(), "ошибка при чтении результатов запроса")
	}

	return packageRecord.toDomain(), nil
}
