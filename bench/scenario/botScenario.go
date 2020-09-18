package scenario

import (
	"context"
	"strconv"
	"sync"

	"github.com/isucon10-qualify/isucon10-qualify/bench/client"
	"github.com/isucon10-qualify/isucon10-qualify/bench/fails"
	"github.com/morikuni/failure"
)

func botScenario(ctx context.Context, c *client.Client) {
	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		q, err := createRandomChairSearchQuery()
		if err != nil {
			fails.Add(err)
		}
		q.Set("perPage", "10")
		chairs, err := c.SearchChairsWithQuery(ctx, q)
		if err != nil {
			code, _ := failure.CodeOf(err)
			if code != fails.ErrBot {
				fails.Add(err)
			}
			return
		}

		for _, chair := range chairs.Chairs {
			wg.Add(1)
			go func(id string) {
				defer wg.Done()
				_, err := c.GetChairDetailFromID(ctx, id)
				code, _ := failure.CodeOf(err)
				if code != fails.ErrBot {
					fails.Add(err)
				}
			}(strconv.FormatInt(chair.ID, 10))
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		q, err := createRandomEstateSearchQuery()
		if err != nil {
			fails.Add(err)
		}
		q.Set("perPage", "10")

		estates, err := c.SearchEstatesWithQuery(ctx, q)
		if err != nil {
			code, _ := failure.CodeOf(err)
			if code != fails.ErrBot {
				fails.Add(err)
			}
			return
		}

		for _, estate := range estates.Estates {
			wg.Add(1)
			go func(id string) {
				defer wg.Done()
				_, err := c.GetEstateDetailFromID(ctx, id)
				code, _ := failure.CodeOf(err)
				if code != fails.ErrBot {
					fails.Add(err)
				}
			}(strconv.FormatInt(estate.ID, 10))
		}
	}()

	wg.Wait()
}
