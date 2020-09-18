package asset

import (
	"encoding/json"
	"reflect"
	"sync"
	"testing"
)

func Test_ParallelStockDecrement(t *testing.T) {
	initialStock := int64(1000000)
	c := Chair{
		stock: initialStock,
	}
	var wg sync.WaitGroup
	start := make(chan struct{})
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			for j := 0; j < 100; j++ {
				c.DecrementStock()
			}
		}()
	}
	close(start)
	wg.Wait()
	got := c.GetStock()
	expected := initialStock - 100*100
	if got != expected {
		t.Errorf("unexpected stocks. expected: %v, but got: %v", expected, got)
	}
}

func TestChair_MarshalJSON(t *testing.T) {
	chair := Chair{
		ID:          1,
		Name:        "name",
		Description: "description",
		Thumbnail:   "thumbnail",
		Price:       2,
		Height:      3,
		Width:       4,
		Depth:       5,
		Color:       "color",
		Features:    "features",
		Kind:        "kind",
		stock:       6,
		popularity:  7,
	}
	b, err := json.Marshal(chair)
	if err != nil {
		t.Fatal("failed to marshal json:", err)
	}
	var got Chair
	if err := json.Unmarshal(b, &got); err != nil {
		t.Fatal("failed to unmarshal json:", err)
	}
	if !reflect.DeepEqual(chair, got) {
		t.Errorf("unexpected chair. expected: %+v, but got: %+v", chair, got)
	}
}
