package handler

import (
	"money-manager/usecase"
	"net/http"

	"github.com/labstack/echo"
)

type UserHandler struct {
	PaymentUsecase usecase.IPaymentUsecase
	// NotifyUsecase  usecase.NotifyUsecase
}

type requestRegisterPayment struct {
	Price int    `json:"price"`
	Name  string `json:"name"`
}

func NewUserHandler(paymentUsecase usecase.IPaymentUsecase) *UserHandler {
	return &UserHandler{
		PaymentUsecase: paymentUsecase,
	}
}

func (uh *UserHandler) RegisterPayment() echo.HandlerFunc {
	return func(c echo.Context) error {
		var rrp requestRegisterPayment
		if err := c.Bind(&rrp); err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		if err := uh.PaymentUsecase.CreatePayment(rrp.Price, rrp.Name); err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, SuccessResult{Status: true, Data: "hoge"})
	}
}

func (uh *UserHandler) GetPaymentResult() echo.HandlerFunc {
	return func(c echo.Context) error {
		result, err := uh.PaymentUsecase.GetPaymentResult()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		return c.JSON(http.StatusOK, SuccessResult{Status: true, Data: result})
	}
}
