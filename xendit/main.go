package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/xendit/xendit-go"
	"github.com/xendit/xendit-go/invoice"
)

var LatestCallback string = ""

func main() {
	// Echo instance
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!\n")
	})

	e.GET("/callback_invoice", func(c echo.Context) error {
		return c.String(http.StatusOK, "Latest Callback: "+LatestCallback)
	})

	e.POST("/callback_invoice", func(c echo.Context) error {
		var json map[string]interface{} = map[string]interface{}{}
		if err := c.Bind(&json); err != nil {
			return err
		}
		LatestCallback = fmt.Sprintf("%s", json)
		return c.JSON(200, json)
	})

	e.GET("/dummy", func(c echo.Context) error {
		redirURL := postDummyXenditInvoice()
		return c.String(200, redirURL)
	})

	e.Logger.Fatal(e.Start(":1323"))
}

func postDummyXenditInvoice() string {

	godotenv.Load()
	xendit.Opt.SecretKey = os.Getenv("XENDIT_API")

	item := xendit.InvoiceItem{
		Name:     "Air Conditioner",
		Quantity: 1,
		Price:    100000,
		Category: "Electronic",
		Url:      "https://yourcompany.com/example_item",
	}

	fee := xendit.InvoiceFee{
		Type:  "ADMIN",
		Value: 5000,
	}

	data := invoice.CreateParams{
		ExternalID:         "terserah_kita_" + uuid.NewString(),
		Amount:             50000,
		Description:        "Invoice Demo #123",
		InvoiceDuration:    86400,
		SuccessRedirectURL: "https://www.google.com",
		FailureRedirectURL: "https://www.google.com",
		Currency:           "IDR",
		Items:              []xendit.InvoiceItem{item},
		Fees:               []xendit.InvoiceFee{fee},
	}

	resp, err := invoice.Create(&data)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Invoice ID", resp.ID)
	fmt.Println("Invoice URL", resp.InvoiceURL)
	return resp.InvoiceURL
	// fmt.Printf("created invoice: %+v\n", resp)
}
