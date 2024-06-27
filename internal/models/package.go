package models

type PackageId int

type PackageName string

type MaxWeight int

type Price int

type Package struct {
	PackageId PackageId   // id типа упаковки
	Name      PackageName // название упаковки
	MaxWeight MaxWeight   // максимальный вес заказа, который можно кпаковать данным способом
	Price     Price       // цена упаковки
}
