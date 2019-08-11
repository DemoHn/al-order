package app

import (
	"fmt"

	"log"

	echo "github.com/labstack/echo"
)

// App - App instance
type App struct {
	*echo.Echo
}

// New - new echo application
func New() *App {
	e := echo.New()
	bindRoutes(e)

	return &App{
		Echo: e,
	}
}

// Start - start server
func (a *App) Start(host string, port int) error {
	addr := fmt.Sprintf("%s:%d", host, port)
	log.Printf("going to listen on address: %s", addr)

	return a.Echo.Start(addr)
}

func bindRoutes(e *echo.Echo) {
	e.POST("/orders", PlaceOrderAPI)
	e.PATCH("/orders/:id", TakeOrderAPI)
	e.GET("/orders", ListOrderAPI)
}
