package client

import (
	"context"
	"strconv"

	"golang.org/x/sync/errgroup"

	"github.com/isucon10-qualify/isucon10-qualify/bench/asset"
)

func (c *Client) AccessTopPage(ctx context.Context) (*ChairsResponse, *EstatesResponse, error) {
	eg, childCtx := errgroup.WithContext(ctx)

	var (
		chairs  *ChairsResponse
		estates *EstatesResponse
	)

	eg.Go(func() (err error) {
		chairs, err = c.GetLowPricedChair(childCtx)
		if err != nil {
			return err
		}
		return nil
	})

	eg.Go(func() (err error) {
		estates, err = c.GetLowPricedEstate(childCtx)
		if err != nil {
			return err
		}
		return nil
	})

	if err := eg.Wait(); err != nil {
		return nil, nil, err
	}

	return chairs, estates, nil
}

func (c *Client) AccessChairDetailPage(ctx context.Context, id int64) (*asset.Chair, *EstatesResponse, error) {
	eg, childCtx := errgroup.WithContext(ctx)

	var (
		chair   *asset.Chair
		estates *EstatesResponse
	)

	eg.Go(func() (err error) {
		chair, err = c.GetChairDetailFromID(childCtx, strconv.FormatInt(id, 10))
		if err != nil {
			return err
		}
		if chair == nil {
			return nil
		}

		return nil
	})

	eg.Go(func() (err error) {
		estates, err = c.GetRecommendedEstatesFromChair(childCtx, id)
		if err != nil {
			return err
		}

		return nil
	})

	if err := eg.Wait(); err != nil {
		return nil, nil, err
	}

	return chair, estates, nil
}

func (c *Client) AccessEstateDetailPage(ctx context.Context, id int64) (*asset.Estate, error) {
	estate, err := c.GetEstateDetailFromID(ctx, strconv.FormatInt(id, 10))
	if err != nil {
		return nil, err
	}

	return estate, nil
}

func (c *Client) AccessChairSearchPage(ctx context.Context) error {
	_, err := c.GetChairSearchCondition(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) AccessEstateSearchPage(ctx context.Context) error {
	_, err := c.GetEstateSearchCondition(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) AccessEstateNazottePage(ctx context.Context) error {
	return nil
}
