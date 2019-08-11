package app

import (
	"database/sql"
	"os"

	util "github.com/DemoHn/al-order/util"
	"github.com/go-redis/redis"
	uuid "github.com/satori/go.uuid"
)

// global function adapters - used for mock tests
var (
	dSaveNewOrder      func(*sql.DB, string, float32, LocationInfo) error
	dFindOneOrder      func(*sql.DB, string) (*Order, error)
	dUpdateOrderStatus func(*sql.DB, string, OrderStatus) error
	dListOrders        func(*sql.DB, int, int) ([]Order, error)
	gdb                *sql.DB
	gredis             *redis.Client
)

// PlaceOrder - place a new order
func PlaceOrder() Error {

}

// TakeOrder - take an existing order
func TakeOrder() {

}

// ListOrder - list recorded orders
func ListOrder() {

}

// private functions
func newUUID() string {
	return uuid.NewV4().String()
}

func getDB() (*sql.DB, error) {
	// return instant db connection
	if gdb != nil {
		return gdb, nil
	}

	dbSource := os.Getenv("DATABASE_URL")
	return util.OpenDB(dbSource)
}

func getRedis() *redis.Client {
	if gredis != nil {
		return gredis
	}

	return redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: "",
		DB:       0,
	})
}

// init function
func init() {
	dSaveNewOrder = SaveNewOrder
	dFindOneOrder = FindOneOrder
	dUpdateOrderStatus = UpdateOrderStatus
	dListOrders = ListOrders
}
