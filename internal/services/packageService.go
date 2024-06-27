package services

import (
	"context"
	"github.com/pkg/errors"
	"gitlab.ozon.dev/go/classroom-13/students/Homework/Homework-1/internal/models"
	"gitlab.ozon.dev/go/classroom-13/students/Homework/Homework-1/internal/storage"
	"gitlab.ozon.dev/go/classroom-13/students/Homework/Homework-1/internal/storage/transactor"
)

type PackageService struct {
	Storage   *storage.PackageStorage
	TxManager *transactor.TransactionManager
}

func NewPackageService(storage *storage.PackageStorage, txManager *transactor.TransactionManager) *PackageService {
	return &PackageService{
		Storage:   storage,
		TxManager: txManager,
	}
}

func (p *PackageService) GetPackageByName(ctx context.Context, packageName string) (models.Package, error) {

	Package, err := p.Storage.GetPackageByName(ctx, models.PackageName(packageName))

	if err != nil {
		return models.Package{}, errors.Wrap(err, "не удалось получить тип упаковки")
	}

	return Package, err
}
