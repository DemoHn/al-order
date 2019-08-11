package app

import (
	"fmt"
	"strconv"

	valid "github.com/asaskevich/govalidator"
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
	locInput := new(LocationInput)
	if err := c.Bind(locInput); err != nil {
		return err
	}
	// validate input
	if err := validateLocationInput(locInput); err != nil {
		return ErrInputValidation(err)
	}

	order, err := sPlaceOrder(*locInput)
	if err != nil {
		return err
	}

	return c.JSON(200, order)
}

// StatusChangeInput -
type StatusChangeInput struct {
	Status string `json:"status"`
}

// TakeOrderAPI - PATCH /orders/:id
func TakeOrderAPI(c echo.Context) error {
	orderID := c.Param("id")
	statusInput := new(StatusChangeInput)
	if err := c.Bind(statusInput); err != nil {
		return err
	}

	if statusInput.Status != "TAKEN" {
		return ErrInputValidation(fmt.Errorf("invalid status action: %s", statusInput.Status))
	}

	if err := sTakeOrder(orderID); err != nil {
		return err
	}
	// finally return data
	return c.JSON(200, map[string]string{
		"status": "SUCCESS",
	})
}

// ListOrderAPI - GET /orders?page=:page&limit=:limit
func ListOrderAPI(c echo.Context) error {
	limitStr := c.QueryParam("limit")
	pageStr := c.QueryParam("page")

	ptrLimit := new(int)
	ptrPage := new(int)

	if limitStr == "" {
		ptrLimit = nil
	} else {
		if !valid.IsInt(limitStr) {
			return ErrInputValidation(fmt.Errorf("invalid `limit` value (must be numeric)"))
		}
		*ptrLimit, _ = strconv.Atoi(limitStr)
	}
	if pageStr == "" {
		ptrPage = nil
	} else {
		if !valid.IsInt(pageStr) {
			return ErrInputValidation(fmt.Errorf("invalid `page` value (must be numeric)"))
		}
		*ptrPage, _ = strconv.Atoi(pageStr)
	}

	orders, err := sListOrder(ptrPage, ptrLimit)
	if err != nil {
		return err
	}

	return c.JSON(200, orders)
}

// private functions

// validateLocationInput - use manual validation to validate
// if locationInput is valid.
// If not, it will return detailed error
func validateLocationInput(locInput *LocationInput) error {
	// validate key
	if len(locInput.Origin) < 2 {
		return fmt.Errorf("`origin` has no enough coordinate data (at least 2)")
	}
	if !valid.IsFloat(locInput.Origin[0]) || !valid.IsFloat(locInput.Origin[1]) {
		return fmt.Errorf("invalid numeric string in `origin`")
	}
	if len(locInput.Destination) < 2 {
		return fmt.Errorf("`destination` has no enough coordinate data (at least 2)")
	}
	if !valid.IsFloat(locInput.Destination[0]) || !valid.IsFloat(locInput.Destination[1]) {
		return fmt.Errorf("invalid numeric string in `destination`")
	}
	return nil
}

func init() {
	sPlaceOrder = PlaceOrder
	sListOrder = ListOrder
	sTakeOrder = TakeOrder
}
