package app

import (
	"database/sql"
	"fmt"
	"reflect"
	"testing"
)

func TestPlaceOrder(t *testing.T) {
	var expOrder = Order{
		Distance:   222.2,
		ID:         "0011-2233",
		sequenceID: 1,
		Location:   LocationInfo{},
		Status:     Unassigned,
	}

	var cursor = 0
	// mock GetRouteDistance()
	gGetRouteDistance = func(loc LocationInfo) (bool, float32, error) {
		var fixtures = []struct {
			hasRoute bool
			distance float32
			err      error
		}{
			{
				hasRoute: true,
				distance: 20.2,
				err:      nil,
			},
			{
				hasRoute: false,
				distance: 0.0,
				err:      nil,
			},
			{
				hasRoute: true,
				distance: 0.0,
				err:      fmt.Errorf("new error"),
			},
		}

		var x = fixtures[cursor]
		return x.hasRoute, x.distance, x.err
	}

	// mock SaveNewOrder()
	dSaveNewOrder = func(*sql.DB, string, float32, LocationInfo) error {
		return nil
	}

	// mock FindOneOrder()
	dFindOneOrder = func(*sql.DB, string) (*Order, error) {
		return &expOrder, nil
	}

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
		{
			name: "normal return value",
			args: args{
				input: LocationInput{
					Origin:      []string{"20.1", "120.2"},
					Destination: []string{"20.1", "120.3"},
				},
			},
			want: &OrderResponse{
				ID:       "0011-2233",
				Distance: 222,
				Status:   Unassigned,
			},
			want1: nil,
		},
		{
			name: "no route (no result)",
			args: args{
				input: LocationInput{
					Origin:      []string{"20.1", "120.2"},
					Destination: []string{"20.1", "120.3"},
				},
			},
			want:  nil,
			want1: ErrNoRoute(),
		},
		{
			name: "google map service error",
			args: args{
				input: LocationInput{
					Origin:      []string{"20.1", "120.2"},
					Destination: []string{"20.1", "120.3"},
				},
			},
			want:  nil,
			want1: ErrGoogleMapService(fmt.Errorf("new error")),
		},
	}
	for idx, tt := range tests {
		cursor = idx

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
	// mock FindOneOrder()
	var cursor = 0
	dFindOneOrder = func(db *sql.DB, orderID string) (*Order, error) {
		if cursor == 0 {
			return &Order{
				Status: Unassigned,
			}, nil
		}

		if cursor == 1 {
			return &Order{
				Status: Taken,
			}, nil
		}
		if cursor == 2 {
			return nil, nil
		}
		return nil, nil
	}

	// mock UpdateOrderStatus() to make all func pass!
	dUpdateOrderStatus = func(db *sql.DB, orderID string, newStatus string) error {
		return nil
	}
	type args struct {
		orderID string
	}
	tests := []struct {
		name string
		args args
		want *Error
	}{
		// TODO: Add test cases.
		{
			name: "should return no error",
			args: args{
				orderID: "0000-1111",
			},
			want: nil,
		},
		{
			name: "error: has been taken",
			args: args{
				orderID: "0000-1111",
			},
			want: ErrOrderHasTaken("from db"),
		},
		{
			name: "error: ID not found",
			args: args{
				orderID: "0000-1111",
			},
			want: ErrIDNotFound(fmt.Sprintf("order id: %s not found", "0000-1111")),
		},
	}
	for idx, tt := range tests {
		cursor = idx
		t.Run(tt.name, func(t *testing.T) {
			if got := TakeOrder(tt.args.orderID); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TakeOrder() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestListOrder(t *testing.T) {
	// mock ListOrders() function
	var records = [][]int{}

	var expOrder = Order{
		Distance:   222.2,
		ID:         "0011-2233",
		sequenceID: 1,
		Location:   LocationInfo{},
		Status:     Unassigned,
	}

	var expOrderResp = OrderResponse{
		Distance: 222,
		ID:       "0011-2233",
		Status:   Unassigned,
	}

	dListOrders = func(db *sql.DB, limit int, offset int) ([]Order, error) {
		records = append(records, []int{limit, offset})
		return []Order{expOrder}, nil
	}

	var rLimit20 = 20
	var rPage3 = 3

	type args struct {
		page  *int
		limit *int
	}
	tests := []struct {
		name      string
		args      args
		mockInput []int
		want      []OrderResponse
		want1     *Error
	}{
		// TODO: Add test cases.
		{
			name: "both limit and offset is null",
			args: args{
				page:  nil,
				limit: nil,
			},
			mockInput: []int{10, 0},
			want:      []OrderResponse{expOrderResp},
			want1:     nil,
		},
		{
			name: "assign limit number",
			args: args{
				page:  nil,
				limit: &rLimit20,
			},
			mockInput: []int{20, 0},
			want:      []OrderResponse{expOrderResp},
			want1:     nil,
		},
		{
			name: "assign limit number & page number",
			args: args{
				page:  &rPage3,
				limit: &rLimit20,
			},
			mockInput: []int{20, 40},
			want:      []OrderResponse{expOrderResp},
			want1:     nil,
		},
	}
	for idx, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := ListOrder(tt.args.page, tt.args.limit)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ListOrder() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("ListOrder() got1 = %v, want %v", got1, tt.want1)
			}
			if !reflect.DeepEqual(records[idx], tt.mockInput) {
				t.Errorf("ListOrder() mock func error: assertInput = %v, record = %v", tt.mockInput, records[idx])
			}
		})
	}
}
