package app

import (
	"fmt"

	echo "github.com/labstack/echo"
)

// global function adapter
var (
	sPlaceOrder func(LocationInput) (*OrderResponse, *Error)
	sTakeOrder  func(string) *Error
	sListOrder  func(*int, *int) ([]OrderResponse, *Error)
)

// PlaceOrderAPI - POST /orders
func PlaceOrderAPI(c echo.Context) error {
	fmt.Println("DIU")
	return nil
}

// TakeOrderAPI - PATCH /orders/:id
func TakeOrderAPI(c echo.Context) error {
	return nil
}

// ListOrderAPI - GET /orders?page=:page&limit=:limit
func ListOrderAPI(c echo.Context) error {
	return nil
}

func init() {
	sPlaceOrder = PlaceOrder
	sListOrder = ListOrder
	sTakeOrder = TakeOrder
}
