package models

import "time"

type OrderId int

type UserId int

type TimeEnd time.Time

type IsGiven bool

type TimeGiven time.Time

type IsReturned bool

type FirstPrice int

type OrderWeight int

type ResultPrice int

type Order struct {
	OrderId     OrderId     // id заказа
	UserId      UserId      // id получателя
	TimeEnd     TimeEnd     // срок хранения заказа
	IsGiven     IsGiven     // флаг - выдан ли заказ
	TimeGiven   TimeGiven   // время выдачи заказа
	IsReturned  IsReturned  // флаг - возвращён ли заказ
	FirstPrice  FirstPrice  // изначальная цена заказа
	Weight      OrderWeight // вес заказа
	PackageName PackageName // тип упаковки
	ResultPrice ResultPrice // окончательная цена заказа
}
