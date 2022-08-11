package main

import (
	"context"
	"fmt"
	"os"

	"money-manager/handler"
	"money-manager/infra"
	"money-manager/usecase"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	log "github.com/labstack/gommon/log"
)

func main() {
	// DI
	ctx := context.Background()
	googleApiHttpClient := infra.NewGoogleApiHttpClient(ctx)
	paymentRepository := infra.NewPaymentRepository(googleApiHttpClient)
	paymentUsecase := usecase.NewPaymentUsecase(paymentRepository)
	userHandler := handler.NewUserHandler(paymentUsecase)

	// Routing
	e := echo.New()
	e.Debug = true
	e.Logger.SetLevel(log.DEBUG)

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/v1/user/payment", userHandler.GetPaymentResult())
	e.POST("/v1/user/payment/register", userHandler.RegisterPayment())

	// e.HTTPErrorHandler = errorResponseHandler

	port := getEnv("SERVER_APP_PORT", "3000")
	e.Logger.Infof("starting money-manager server on : %s ...", port)
	serverPort := fmt.Sprintf(":%s", port)
	e.Logger.Fatal(e.Start(serverPort))
}

func getEnv(key string, defaultValue string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return defaultValue
}
