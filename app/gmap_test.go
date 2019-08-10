package app

import (
	"testing"
)

func TestGetRouteDistance(t *testing.T) {

	GetRouteDistance(LocationInfo{
		StartLat: "-30.1",
		StartLng: "120.1",
		EndLat:   "30.12",
		EndLng:   "120.1",
	})
	/*
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
		}
		for _, tt := range tests {
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
		}*/
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
