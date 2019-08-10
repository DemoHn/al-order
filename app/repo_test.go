package app

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"
	"testing"

	// import mysql driver
	"github.com/DemoHn/al-order/util"
	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func TestMain(m *testing.M) {
	if err := setup(); err != nil {
		log.Fatal("setup testcase framework error:", err)
		os.Exit(-1)
		return
	}
	code := m.Run()
	if err := teardown(); err != nil {
		log.Fatal("teardown testcase framework error:", err)
		os.Exit(-1)
		return
	}
	os.Exit(code)
}

// setup & teardown testcases globally
func setup() error {
	m, err := util.ParseEnvFromFile("../.env")
	if err != nil {
		return err
	}
	dbURL := m["DATABASE_URL"]
	fmt.Println(dbURL)
	if err := util.ExecMigration(dbURL, "../sql/up.sql"); err != nil {
		return err
	}

	dbURLs := strings.Split(dbURL, "://")
	var e error
	if db, e = sql.Open(dbURLs[0], dbURLs[1]); e != nil {
		return e
	}
	return nil
}
func teardown() error {
	m, err := util.ParseEnvFromFile("../.env")
	if err != nil {
		return err
	}

	if err := util.ExecMigration(m["DATABASE_URL"], "../sql/down.sql"); err != nil {
		return err
	}
	if err := db.Close(); err != nil {
		return err
	}
	return nil
}

func TestSaveNewOrder(t *testing.T) {
	type args struct {
		id       string
		distance float32
		location LocationInfo
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "insert data normally",
			args: args{
				id:       "0000",
				distance: 12.5,
				location: LocationInfo{},
			},
			wantErr: false,
		},
		{
			name: "insert duplicate data",
			args: args{
				id:       "0000",
				distance: 13.5,
				location: LocationInfo{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := SaveNewOrder(db, tt.args.id, tt.args.distance, tt.args.location); (err != nil) != tt.wantErr {
				t.Errorf("SaveNewOrder() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFindOneOrder(t *testing.T) {
	type args struct {
		id string
	}
	tests := []struct {
		name    string
		args    args
		want    *Order
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "find one order",
			args: args{
				id: "0000",
			},
			want: &Order{
				ID:         "0000",
				sequenceID: 1,
				Location:   LocationInfo{},
				Distance:   12.5,
				Status:     Unassigned,
			},
			wantErr: false,
		},
		{
			name: "find no record (and won't return error)",
			args: args{
				id: "1234567_not_exists",
			},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FindOneOrder(db, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("FindOneOrder() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FindOneOrder() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUpdateOrderStatus(t *testing.T) {
	type args struct {
		id        string
		newStatus OrderStatus
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "update status -> TAKEN",
			args: args{
				id:        "0000",
				newStatus: Taken,
			},
			wantErr: false,
		},
		{
			name: "won't return error even id not exists",
			args: args{
				id:        "not_exists",
				newStatus: Taken,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := UpdateOrderStatus(db, tt.args.id, tt.args.newStatus); (err != nil) != tt.wantErr {
				t.Errorf("UpdateOrderStatus() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_encodeLocationInfo(t *testing.T) {
	type args struct {
		info LocationInfo
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "if empty",
			args: args{
				info: LocationInfo{},
			},
			want: "{\"StartLat\":0,\"StartLng\":0,\"EndLat\":0,\"EndLng\":0}",
		},
		{
			name: "normal data",
			args: args{
				info: LocationInfo{
					StartLat: -45.3,
					StartLng: 123.6,
					EndLat:   12,
					EndLng:   12.234,
				},
			},
			want: "{\"StartLat\":-45.3,\"StartLng\":123.6,\"EndLat\":12,\"EndLng\":12.234}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := encodeLocationInfo(tt.args.info); got != tt.want {
				t.Errorf("encodeLocationInfo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_decodeLocationInfo(t *testing.T) {
	type args struct {
		data string
	}
	tests := []struct {
		name string
		args args
		want LocationInfo
	}{
		// TODO: Add test cases.
		{
			name: "if empty",
			args: args{
				data: "{\"StartLat\":0,\"StartLng\":0,\"EndLat\":0,\"EndLng\":0}",
			},
			want: LocationInfo{},
		},
		{
			name: "if has data",
			args: args{
				data: "{\"StartLat\":-45.3,\"StartLng\":123.6,\"EndLat\":12,\"EndLng\":12.234}",
			},
			want: LocationInfo{
				StartLat: -45.3,
				StartLng: 123.6,
				EndLat:   12,
				EndLng:   12.234,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := decodeLocationInfo(tt.args.data); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("decodeLocationInfo() = %v, want %v", got, tt.want)
			}
		})
	}
}
