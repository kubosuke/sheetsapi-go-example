package infra

import "time"

type PaymentDto struct {
	Date     time.Time
	Name     string
	Price    int
	Claiment string
}
