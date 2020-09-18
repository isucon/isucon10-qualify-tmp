package scenario

import (
	"context"
	"log"

	"github.com/isucon10-qualify/isucon10-qualify/bench/client"
	"github.com/isucon10-qualify/isucon10-qualify/bench/fails"
	"github.com/isucon10-qualify/isucon10-qualify/bench/parameter"
	"github.com/morikuni/failure"
)

func Initialize(ctx context.Context) *client.InitializeResponse {
	// Initializeにはタイムアウトを設定
	// レギュレーションにある時間を設定する
	ctx, cancel := context.WithTimeout(ctx, parameter.InitializeTimeout)
	defer cancel()

	res, err := initialize(ctx)
	if err != nil {
		if ctx.Err() != nil {
			err = failure.New(fails.ErrCritical, failure.Message("POST /initialize: リクエストがタイムアウトしました"))
			fails.Add(err)
		} else {
			fails.Add(err)
		}
	}
	return res
}

func Validation(ctx context.Context) {
	cancelCtx, cancel := context.WithTimeout(ctx, parameter.LoadTimeout)
	defer cancel()
	go Load(cancelCtx)

	select {
	case <-fails.Fail():
		log.Println("fail条件を満たしました")
		return
	case <-cancelCtx.Done():
		return
	}
}
