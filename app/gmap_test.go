package app

import (
	"fmt"
	"testing"
)

var cursor int

func init() {
	cursor = 0
	var fixtures = []struct {
		data []byte
		err  error
	}{
		{
			data: []byte("{\"status\":\"OK\", \"routes\": [{\"legs\": [{\"distance\": {\"text\":\"233\", \"value\": 233}}]}]}"),
			err:  nil,
		},
		{
			data: []byte("{\"status\":\"ZERO_RESULTS\", \"routes\": []}"),
			err:  nil,
		},
		{
			data: []byte{},
			err:  fmt.Errorf("new error"),
		},
	}
	// replace default function with mock function
	fetchData = func(string) ([]byte, error) {
		return fixtures[cursor].data, fixtures[cursor].err
	}
}
func TestGetRouteDistance(t *testing.T) {

	type args struct {
		loc LocationInfo
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		want1   float32
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "normal get distance",
			args: args{
				loc: LocationInfo{},
			},
			want:    true,
			want1:   233.0,
			wantErr: false,
		},
		{
			name: "no valid routes",
			args: args{
				loc: LocationInfo{},
			},
			want:    false,
			want1:   0.0,
			wantErr: false,
		},
		{
			name: "encounter error",
			args: args{
				loc: LocationInfo{},
			},
			want:    false,
			want1:   0.0,
			wantErr: true,
		},
	}
	for index, tt := range tests {
		cursor = index
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := GetRouteDistance(tt.args.loc)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetRouteDistance() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetRouteDistance() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("GetRouteDistance() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_formatPoints(t *testing.T) {
	type args struct {
		loc LocationInfo
	}
	tests := []struct {
		name  string
		args  args
		want  string
		want1 string
	}{
		// TODO: Add test cases.
		{
			name: "example format data",
			args: args{
				loc: LocationInfo{
					StartLat: "1.2",
					StartLng: "20.1",
					EndLat:   "1.2",
					EndLng:   "20.2",
				},
			},
			want:  "1.2,20.1",
			want1: "1.2,20.2",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := formatPoints(tt.args.loc)
			if got != tt.want {
				t.Errorf("formatPoints() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("formatPoints() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
