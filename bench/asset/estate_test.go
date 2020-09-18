package asset

import (
	"testing"
)

func TestEstate_ToCSV(t *testing.T) {
	type fields struct {
		ID          int64
		Thumbnail   string
		Name        string
		Description string
		Address     string
		Latitude    float64
		Longitude   float64
		DoorHeight  int64
		DoorWidth   int64
		Rent        int64
		Features    string
		popularity  int64
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			fields: fields{
				ID:          1,
				Thumbnail:   "/foo/bar",
				Name:        "chair",
				Description: "説明",
				Address:     "東京都",
				Latitude:    1.2,
				Longitude:   3.45,
				DoorHeight:  10,
				DoorWidth:   20,
				Rent:        30,
				Features:    "a,b,c,d,e",
				popularity:  3,
			},
			want: `1,chair,説明,/foo/bar,東京都,1.2,3.45,30,10,20,"a,b,c,d,e",3` + "\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Estate{
				ID:          tt.fields.ID,
				Thumbnail:   tt.fields.Thumbnail,
				Name:        tt.fields.Name,
				Description: tt.fields.Description,
				Address:     tt.fields.Address,
				Latitude:    tt.fields.Latitude,
				Longitude:   tt.fields.Longitude,
				DoorHeight:  tt.fields.DoorHeight,
				DoorWidth:   tt.fields.DoorWidth,
				Rent:        tt.fields.Rent,
				Features:    tt.fields.Features,
				popularity:  tt.fields.popularity,
			}
			if got := e.ToCSV(); got != tt.want {
				t.Errorf("ToCSV() = %v, want %v", got, tt.want)
			}
		})
	}
}
