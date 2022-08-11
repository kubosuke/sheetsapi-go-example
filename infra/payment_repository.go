package infra

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"golang.org/x/oauth2/jwt"
	"google.golang.org/api/sheets/v4"
)

type paymentRepository struct {
	googleApiHttpClient *http.Client
}

func NewPaymentRepository(googleApiHttpClient *http.Client) *paymentRepository {
	return &paymentRepository{googleApiHttpClient: googleApiHttpClient}
}

func NewGoogleApiHttpClient(ctx context.Context) *http.Client {
	conf := &jwt.Config{
		Email:        EMAIL,
		PrivateKey:   []byte(PRIVATE_KEY),
		PrivateKeyID: PRIVATE_KEY_ID,
		TokenURL:     TOKEN_URL,
		Scopes: []string{
			"https://www.googleapis.com/auth/spreadsheets",
		},
	}
	client := conf.Client(ctx)
	client.Timeout = time.Second * 10
	return client
}

func (p *paymentRepository) Create(sheetName, date, name, claiment string, price int) error {
	srv, err := sheets.New(p.googleApiHttpClient)
	if err != nil {
		return fmt.Errorf("Unable to retrieve Sheets client: %v", err)
	}

	var values [][]any
	values = append(values, []any{date, name, price, claiment})
	if _, err = srv.Spreadsheets.Values.Append(SPREAD_SHEET_ID, sheetName, &sheets.ValueRange{
		Values: values,
	}).ValueInputOption("RAW").InsertDataOption("INSERT_ROWS").Do(); err != nil {
		return fmt.Errorf("Unable to write Sheets: %v", err)
	}
	return nil
}

func (p *paymentRepository) IsSheetExist(sheetName string) (bool, error) {
	srv, err := sheets.New(p.googleApiHttpClient)
	if err != nil {
		return false, fmt.Errorf("Unable to retrieve Sheets client: %v", err)
	}
	sheets, err := srv.Spreadsheets.Get(SPREAD_SHEET_ID).Do()
	if err != nil {
		return false, fmt.Errorf("Unable to retrieve Sheets: %v", err)
	}
	for _, sheet := range sheets.Sheets {
		if sheet.Properties.Title == sheetName {
			return true, nil
		}
	}
	return false, nil
}

func (p *paymentRepository) CreateSheet(sheetName string) error {
	srv, err := sheets.New(p.googleApiHttpClient)
	if err != nil {
		return fmt.Errorf("Unable to retrieve Sheets client: %v", err)
	}
	properties := sheets.SheetProperties{Title: sheetName}
	addSheet := sheets.AddSheetRequest{Properties: &properties}
	request := sheets.Request{AddSheet: &addSheet}
	batchUpdateSpreadsheetRequest := sheets.BatchUpdateSpreadsheetRequest{Requests: []*sheets.Request{&request}}
	if _, err := srv.Spreadsheets.BatchUpdate(SPREAD_SHEET_ID, &batchUpdateSpreadsheetRequest).Do(); err != nil {
		return fmt.Errorf("Unable to Create Sheet: %v", err)
	}
	return nil
}

func (p *paymentRepository) ListPayment(sheetName string) ([]PaymentDto, error) {
	srv, err := sheets.New(p.googleApiHttpClient)
	if err != nil {
		return nil, fmt.Errorf("Unable to retrieve Sheets client: %v", err)
	}
	data, err := srv.Spreadsheets.Values.Get(SPREAD_SHEET_ID, sheetName).Do()
	if err != nil {
		return nil, fmt.Errorf("Unable to retrieve Sheets data: %v", err)
	}

	var payments []PaymentDto
	for _, payment_info := range data.Values {
		date, ok := payment_info[0].(time.Time)
		if !ok {
			log.Printf("unexpected date: %s", date)
			continue
		}
		name, ok := payment_info[1].(string)
		if !ok {
			log.Printf("unexpected name: %s", date)
			continue
		}
		price, ok := payment_info[2].(int)
		if !ok {
			log.Printf("unexpected price: %s", date)
			continue
		}
		claiment, ok := payment_info[3].(string)
		if !ok {
			log.Printf("unexpected claiment: %s", date)
			continue
		}
		payments = append(payments, PaymentDto{
			date,
			name,
			price,
			claiment,
		})
	}
	return payments, nil
}
