package app

import (
	"database/sql"
	"reflect"
	"testing"

	"github.com/go-redis/redis"
)

func TestPlaceOrder(t *testing.T) {
	type args struct {
		input LocationInput
	}
	tests := []struct {
		name  string
		args  args
		want  *OrderResponse
		want1 *Error
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := PlaceOrder(tt.args.input)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PlaceOrder() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("PlaceOrder() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestTakeOrder(t *testing.T) {
	type args struct {
		orderID string
	}
	tests := []struct {
		name string
		args args
		want *Error
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := TakeOrder(tt.args.orderID); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TakeOrder() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestListOrder(t *testing.T) {
	type args struct {
		page  *int
		limit *int
	}
	tests := []struct {
		name  string
		args  args
		want  []OrderResponse
		want1 *Error
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := ListOrder(tt.args.page, tt.args.limit)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ListOrder() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("ListOrder() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_newUUID(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newUUID(); got != tt.want {
				t.Errorf("newUUID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getDB(t *testing.T) {
	tests := []struct {
		name    string
		want    *sql.DB
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getDB()
			if (err != nil) != tt.wantErr {
				t.Errorf("getDB() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getDB() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getRedis(t *testing.T) {
	tests := []struct {
		name string
		want *redis.Client
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getRedis(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getRedis() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_transformOrderResp(t *testing.T) {
	type args struct {
		order *Order
	}
	tests := []struct {
		name string
		args args
		want *OrderResponse
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := transformOrderResp(tt.args.order); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("transformOrderResp() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_addOrderLock(t *testing.T) {
	type args struct {
		client  *redis.Client
		orderID string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := addOrderLock(tt.args.client, tt.args.orderID)
			if (err != nil) != tt.wantErr {
				t.Errorf("addOrderLock() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("addOrderLock() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_deleteOrderLock(t *testing.T) {
	type args struct {
		client  *redis.Client
		orderID string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := deleteOrderLock(tt.args.client, tt.args.orderID); (err != nil) != tt.wantErr {
				t.Errorf("deleteOrderLock() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
