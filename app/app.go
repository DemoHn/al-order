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

// bind routes
func bindRoutes(e *echo.Echo) {
	e.POST("/orders", PlaceOrderAPI)
	e.PATCH("/orders/:id", TakeOrderAPI)
	e.GET("/orders", ListOrderAPI)
}

func bindMiddlewares(e *echo.Echo) {
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if err := next(c); err != nil {
				// wrap errors to export
				var wrapError *Error
				switch e := err.(type) {
				case *Error:
					wrapError = e
				case *echo.HTTPError:
					wrapError = ErrGeneralHTTP(e)
				default:
					wrapError = ErrUnknownFatal(err)
				}

				// print error log
				log.Fatalf("Error(%d): %s\n Detail: %s", wrapError.Code, wrapError.Message, wrapError.Details)

				return c.JSON(wrapError.StatusCode, map[string]string{
					"error": wrapError.Message,
				})
			}
			return nil
		}
	})
}
