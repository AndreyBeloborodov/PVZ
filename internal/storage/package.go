package storage

import (
	"gitlab.ozon.dev/go/classroom-13/students/Homework/Homework-1/internal/models"
)

type PackageRecord struct {
	PackageId int    `db:"id"`
	Name      string `db:"name"`
	MaxWeight int    `db:"max_weight"`
	Price     int    `db:"price"`
}

func (p PackageRecord) toDomain() models.Package {
	return models.Package{
		PackageId: models.PackageId(p.PackageId),
		Name:      models.PackageName(p.Name),
		MaxWeight: models.MaxWeight(p.MaxWeight),
		Price:     models.Price(p.Price),
	}
}

func transformPackage(Package models.Package) PackageRecord {
	return PackageRecord{
		PackageId: int(Package.PackageId),
		Name:      string(Package.Name),
		MaxWeight: int(Package.MaxWeight),
		Price:     int(Package.Price),
	}
}
