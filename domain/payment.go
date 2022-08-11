package domain

import (
	"time"
)

type Payment struct {
	Date     time.Time
	Name     string
	Price    int
	Claiment string
}

type Payments []Payment

func NewPayment(date time.Time, name string, price int, claiment string) Payment {
	return Payment{
		Date:     date,
		Name:     name,
		Price:    price,
		Claiment: claiment,
	}
}

func (p *Payment) GetDateYYYYMM() string {
	return p.Date.Format(YYYY_MM)
}

func (p *Payment) GetDateYYYYMMDD() string {
	return p.Date.Format(YYYY_MM_DD)
}
