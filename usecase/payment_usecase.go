package usecase

import (
	"fmt"
	"money-manager/domain"
	"money-manager/repository"
	"time"
)

type IPaymentUsecase interface {
	CreatePayment(price int, name string) error
	GetPaymentResult() (domain.Payment, error)
}

type paymentUsecase struct {
	paymentRepository repository.PaymentRepository
}

func NewPaymentUsecase(paymentRepository repository.PaymentRepository) *paymentUsecase {
	return &paymentUsecase{
		paymentRepository: paymentRepository,
	}
}

func (pu *paymentUsecase) CreatePayment(price int, name string) error {
	payment := domain.NewPayment(time.Now(), name, price, "foo")

	isSheetExist, err := pu.paymentRepository.IsSheetExist(payment.Date.Format(domain.YYYY_MM))
	if err != nil {
		return fmt.Errorf("failed to get sheet info. error is: %w", err)
	}

	if !isSheetExist {
		if err := pu.paymentRepository.CreateSheet(payment.Date.Format(domain.YYYY_MM)); err != nil {
			return fmt.Errorf("failed to create sheet. error is: %w", err)
		}
	}

	if err := pu.paymentRepository.Create(payment.Date.Format(domain.YYYY_MM), payment.Date.Format(payment.GetDateYYYYMMDD()), payment.Name, payment.Claiment, payment.Price); err != nil {
		return fmt.Errorf("failed to create payment. error is: %w", err)
	}
	return nil
}

func (pu *paymentUsecase) GetPaymentResult() (domain.Payment, error) {
	payments_info, err := pu.paymentRepository.ListPayment(time.Now().Format(domain.YYYY_MM))
	if err != nil {
		return domain.Payment{}, fmt.Errorf("failed to list payment. error is: %w", err)
	}

	var payments domain.Payments
	for _, payment := range payments_info {
		payments = append(payments, domain.NewPayment(payment.Date, payment.Name, payment.Price, payment.Claiment))
	}

	// payments.CalcPayment()

	return domain.Payment{}, nil
}
