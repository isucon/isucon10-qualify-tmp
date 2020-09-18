package scenario

import (
	"context"
	"math/rand"
	"strconv"
	"time"

	"github.com/morikuni/failure"

	"github.com/isucon10-qualify/isucon10-qualify/bench/asset"
	"github.com/isucon10-qualify/isucon10-qualify/bench/client"
	"github.com/isucon10-qualify/isucon10-qualify/bench/fails"
	"github.com/isucon10-qualify/isucon10-qualify/bench/parameter"
)

func chairSearchScenario(ctx context.Context, c *client.Client) error {
	t := time.Now()
	chairs, estates, err := c.AccessTopPage(ctx)
	if err != nil {
		fails.Add(err)
		return failure.New(fails.ErrApplication)
	}

	if err := checkChairsOrderedByPrice(chairs.Chairs, t); err != nil {
		err = failure.Translate(err, fails.ErrApplication, failure.Message("GET /api/chair/low_priced: レスポンスの内容が不正です"))
		fails.Add(err)
		return failure.New(fails.ErrApplication)
	}

	if err := checkEstatesOrderedByRent(estates.Estates); err != nil {
		err = failure.Translate(err, fails.ErrApplication, failure.Message("GET /api/estate/low_priced: レスポンスの内容が不正です"))
		fails.Add(err)
		return failure.New(fails.ErrApplication)
	}

	if time.Since(t) > parameter.ThresholdTimeOfAbandonmentPage {
		return failure.New(fails.ErrTimeout)
	}

	t = time.Now()
	err = c.AccessChairSearchPage(ctx)
	if err != nil {
		fails.Add(err)
		return failure.New(fails.ErrApplication)
	}
	if time.Since(t) > parameter.ThresholdTimeOfAbandonmentPage {
		return failure.New(fails.ErrTimeout)
	}

	// Search Chairs with Query
	var cr *client.ChairsResponse
	for i := 0; i < parameter.NumOfSearchChairInScenario; i++ {
		q, err := createRandomChairSearchQuery()
		if err != nil {
			fails.Add(err)
			return failure.New(fails.ErrApplication)
		}

		t = time.Now()
		_cr, err := c.SearchChairsWithQuery(ctx, q)
		if err != nil {
			fails.Add(err)
			return failure.New(fails.ErrApplication)
		}

		if time.Since(t) > parameter.ThresholdTimeOfAbandonmentPage {
			return failure.New(fails.ErrTimeout)
		}

		if len(_cr.Chairs) == 0 {
			continue
		}

		if err := checkChairsOrderedByPopularity(_cr.Chairs, t); err != nil {
			err = failure.Translate(err, fails.ErrApplication, failure.Message("GET /api/chair/search: レスポンスの内容が不正です"))
			fails.Add(err)
			return failure.New(fails.ErrApplication)
		}

		cr = _cr
		numOfPages := int(_cr.Count) / parameter.PerPageOfChairSearch
		if numOfPages == 0 {
			continue
		}
		if numOfPages > parameter.LimitOfChairSearchPageDepth {
			numOfPages = parameter.LimitOfChairSearchPageDepth
		}

		for j := 0; j < parameter.NumOfCheckChairSearchPaging; j++ {
			q.Set("page", strconv.Itoa(rand.Intn(numOfPages)))

			t := time.Now()
			_cr, err := c.SearchChairsWithQuery(ctx, q)
			if err != nil {
				fails.Add(err)
				return failure.New(fails.ErrApplication)
			}

			if time.Since(t) > parameter.ThresholdTimeOfAbandonmentPage {
				return failure.New(fails.ErrTimeout)
			}

			if len(_cr.Chairs) == 0 {
				fails.Add(err)
				return failure.New(fails.ErrApplication)
			}

			if err := checkChairsOrderedByPopularity(_cr.Chairs, t); err != nil {
				err = failure.Translate(err, fails.ErrApplication, failure.Message("GET /api/chair/search: レスポンスの内容が不正です"))
				fails.Add(err)
				return failure.New(fails.ErrApplication)
			}

			cr = _cr
			numOfPages = int(cr.Count) / parameter.PerPageOfChairSearch
			if numOfPages == 0 {
				break
			}
			if numOfPages > parameter.LimitOfChairSearchPageDepth {
				numOfPages = parameter.LimitOfChairSearchPageDepth
			}
		}
	}

	if cr == nil || len(cr.Chairs) == 0 {
		return nil
	}

	// Get detail of Chair
	var targetID int64 = -1
	var chair *asset.Chair
	var er *client.EstatesResponse
	for i := 0; i < parameter.NumOfCheckChairDetailPage; i++ {
		randomPosition := rand.Intn(len(cr.Chairs))
		targetID = cr.Chairs[randomPosition].ID
		t = time.Now()
		chair, er, err = c.AccessChairDetailPage(ctx, targetID)

		if err != nil {
			fails.Add(err)
			return failure.New(fails.ErrApplication)
		}

		if time.Since(t) > parameter.ThresholdTimeOfAbandonmentPage {
			return failure.New(fails.ErrTimeout)
		}

		if chair == nil || len(er.Estates) == 0 {
			return nil
		}

		if err := checkChairEqualToAsset(chair); err != nil {
			err = failure.Translate(err, fails.ErrApplication, failure.Message("GET /api/chair/:id: レスポンスの内容が不正です"))
			fails.Add(err)
			return failure.New(fails.ErrApplication)
		}

		if err := checkRecommendedEstates(er.Estates, chair); err != nil {
			err = failure.Translate(err, fails.ErrApplication, failure.Message("GET /api/recommended_estate/:id: レスポンスの内容が不正です"))
			fails.Add(err)
			return failure.New(fails.ErrApplication)
		}
	}

	if targetID == -1 {
		return nil
	}

	// Buy Chair
	err = c.BuyChair(ctx, strconv.FormatInt(targetID, 10))
	if err != nil {
		if _chair, err := asset.GetChairFromID(targetID); err != nil || _chair.GetStock() > 0 {
			fails.Add(err)
			return failure.New(fails.ErrApplication)
		}
	}

	// Get detail of Estate
	targetID = -1
	for i := 0; i < parameter.NumOfCheckEstateDetailPage; i++ {
		randomPosition := rand.Intn(len(er.Estates))
		targetID = er.Estates[randomPosition].ID
		t = time.Now()
		e, err := c.AccessEstateDetailPage(ctx, targetID)
		if err != nil {
			fails.Add(err)
			return failure.New(fails.ErrApplication)
		}

		if time.Since(t) > parameter.ThresholdTimeOfAbandonmentPage {
			return failure.New(fails.ErrTimeout)
		}

		if err := checkEstateEqualToAsset(e); err != nil {
			err = failure.Translate(err, fails.ErrApplication, failure.Message("GET /api/estate/:id: レスポンスの内容が不正です"))
			fails.Add(err)
			return failure.New(fails.ErrApplication)
		}
	}

	if targetID == -1 {
		return nil
	}

	// Request docs of Estate
	err = c.RequestEstateDocument(ctx, strconv.FormatInt(targetID, 10))

	if err != nil {
		fails.Add(err)
		return failure.New(fails.ErrApplication)
	}

	return nil
}
