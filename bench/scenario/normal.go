package scenario

import (
	"context"

	"github.com/isucon10-qualify/isucon10-qualify/bench/client"
)

func initialize(ctx context.Context) (*client.InitializeResponse, error) {
	c := client.NewClientForInitialize()
	res, err := c.Initialize(ctx)
	if err != nil {
		return res, err
	}
	return res, nil
}
