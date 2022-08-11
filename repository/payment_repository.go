package repository

import "money-manager/infra"

type PaymentRepository interface {
	Create(sheetName, date, name, claiment string, price int) error
	IsSheetExist(sheetName string) (bool, error)
	CreateSheet(sheetName string) error
	ListPayment(sheetName string) ([]infra.PaymentDto, error)
}
