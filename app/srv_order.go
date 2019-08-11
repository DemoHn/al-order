package app

import (
	"database/sql"
	"fmt"
	"os"
	"time"

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
	rAddOrderLock      func(*redis.Client, string) (bool, error)
	rDelOrderLock      func(*redis.Client, string) error
	gdb                *sql.DB
	gredis             *redis.Client
)

// types

// LocationInput - location input request schema
type LocationInput struct {
	Origin      []string `json:"origin"`
	Destination []string `json:"destination"`
}

// OrderResponse - show the response of order
type OrderResponse struct {
	ID       string `json:"id"`
	Distance int    `json:"distance"`
	Status   string `json:"status"`
}

// PlaceOrder - place a new order
func PlaceOrder(input LocationInput) (*OrderResponse, *Error) {
	var loc = LocationInfo{
		StartLat: input.Origin[0],
		StartLng: input.Origin[1],
		EndLat:   input.Destination[0],
		EndLng:   input.Destination[1],
	}
	db, err := getDB()
	if err != nil {
		return nil, ErrDBFatal(err)
	}
	orderID := newUUID()

	// 01. get distance data from GoogleMap direction API
	hasRoute, distance, err := GetRouteDistance(loc)
	if err != nil {
		return nil, ErrGoogleMapService(err)
	}
	if hasRoute == false {
		return nil, ErrNoRoute()
	}

	// 02. write down data
	if err := dSaveNewOrder(db, orderID, distance, loc); err != nil {
		return nil, ErrDBFatal(err)
	}

	// 03. find & return the final data
	order, err := dFindOneOrder(db, orderID)
	if err != nil {
		return nil, ErrDBFatal(err)
	}

	return transformOrderResp(order), nil
}

// TakeOrder - take an existing order
func TakeOrder(orderID string) *Error {
	r := getRedis()
	// 01. add lock to avoid unexpected race-condition - that
	// both two users take the order at the same time
	lockResult, err := rAddOrderLock(r, orderID)
	if err != nil {
		return ErrRedisFatal(err)
	}
	defer rDelOrderLock(r, orderID)
	if lockResult == false {
		return ErrOrderHasTaken()
	}

	// 02. find if orderID exists
	db, err := getDB()
	if err != nil {
		return ErrDBFatal(err)
	}
	order, err := dFindOneOrder(db, orderID)
	if err != nil {
		return ErrDBFatal(err)
	}
	if order == nil {
		return ErrIDNotFound(fmt.Sprintf("order id: %s not found", orderID))
	}
	// 03. check if orderID status is still unassigned
	if order.Status != Unassigned {
		return ErrOrderHasTaken()
	}
	// 04. so it's time to update the status!
	if err := dUpdateOrderStatus(db, orderID, Taken); err != nil {
		return ErrDBFatal(err)
	}
	return nil
}

// ListOrder - list recorded orders
func ListOrder(page *int, limit *int) ([]OrderResponse, *Error) {
	var rPage, rLimit int
	// get page, limit
	if page == nil {
		rPage = 1
	} else {
		rPage = *page
	}

	if limit == nil {
		rLimit = 10
	} else {
		rLimit = *limit
	}

	var rOffset = (rPage - 1) * rLimit

	// 01. list orders
	db, err := getDB()
	if err != nil {
		return nil, ErrDBFatal(err)
	}
	orders, err := dListOrders(db, rLimit, rOffset)
	if err != nil {
		return nil, ErrDBFatal(err)
	}

	// 02. transform orders
	var tOrders = make([]OrderResponse, len(orders))
	for index, order := range orders {
		tOrders[index] = *transformOrderResp(&order)
	}
	return tOrders, nil
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

func transformOrderResp(order *Order) *OrderResponse {
	return &OrderResponse{
		ID:       order.ID,
		Status:   order.Status,
		Distance: int(order.Distance),
	}
}

// redis helper
func addOrderLock(client *redis.Client, orderID string) (bool, error) {
	key := fmt.Sprintf("order_lock:%s", orderID)
	return client.SetNX(key, "1", 10*time.Second).Result()
}

func deleteOrderLock(client *redis.Client, orderID string) error {
	key := fmt.Sprintf("order_lock:%s", orderID)
	if _, err := client.Del(key).Result(); err != nil {
		return err
	}
	return nil
}

// init function & adapters
func init() {
	dSaveNewOrder = SaveNewOrder
	dFindOneOrder = FindOneOrder
	dUpdateOrderStatus = UpdateOrderStatus
	dListOrders = ListOrders
	rAddOrderLock = addOrderLock
	rDelOrderLock = deleteOrderLock
}
