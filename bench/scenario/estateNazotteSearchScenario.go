package scenario

import (
	"context"
	"math/rand"
	"sort"
	"strconv"
	"time"

	"github.com/morikuni/failure"

	"github.com/isucon10-qualify/isucon10-qualify/bench/asset"
	"github.com/isucon10-qualify/isucon10-qualify/bench/client"
	"github.com/isucon10-qualify/isucon10-qualify/bench/fails"
	"github.com/isucon10-qualify/isucon10-qualify/bench/parameter"
)

type point struct {
	Latitude  float64
	Longitude float64
}

var famousPlaces []point = []point{
	{Latitude: 34.81667, Longitude: 137.4},
	{Latitude: 34.4833, Longitude: 136.84186},
	{Latitude: 36.65, Longitude: 138.31667},
	{Latitude: 34.9, Longitude: 137.5},
	{Latitude: 35.06667, Longitude: 135.21667},
	{Latitude: 36, Longitude: 139.55722},
	{Latitude: 36.53333, Longitude: 136.61667},
	{Latitude: 36.75965, Longitude: 137.36215},
	{Latitude: 35, Longitude: 136.51667},
	{Latitude: 33.4425, Longitude: 129.96972},
	{Latitude: 35.30889, Longitude: 139.55028},
	{Latitude: 34.25, Longitude: 135.31667},
	{Latitude: 35.82756, Longitude: 137.95378},
	{Latitude: 33.3213, Longitude: 130.94098},
	{Latitude: 36.24624, Longitude: 139.07204},
	{Latitude: 36.33011, Longitude: 138.89585},
	{Latitude: 35.815, Longitude: 139.6853},
	{Latitude: 39.46667, Longitude: 141.95},
	{Latitude: 37.56667, Longitude: 140.11667},
	{Latitude: 43.82634, Longitude: 144.09638},
	{Latitude: 44.35056, Longitude: 142.45778},
	{Latitude: 41.77583, Longitude: 140.73667},
	{Latitude: 35.48199, Longitude: 137.02166},
}

const (
	rangeDiffLatitude  = 3
	rangeDiffLongitude = 3
	rangeMaxWidth      = 1.0
	rangeMinWidth      = 0.1
	rangeMaxHeight     = 1.0
	rangeMinHeight     = 0.1
	numOfMaxPoints     = 20
	numOfMinPoints     = 10
)

func createRandomConvexhull() []point {
	famousPlace := famousPlaces[rand.Intn(len(famousPlaces))]

	width := rand.Float64()*(rangeMaxWidth-rangeMinWidth) + rangeMinWidth
	height := rand.Float64()*(rangeMaxHeight-rangeMinHeight) + rangeMinHeight
	center := point{
		Latitude:  famousPlace.Latitude + (rand.Float64()-0.5)*rangeDiffLatitude,
		Longitude: famousPlace.Longitude + (rand.Float64()-0.5)*rangeDiffLongitude,
	}

	pointCounts := rand.Intn(numOfMaxPoints-numOfMinPoints) + numOfMinPoints

	points := []point{}

	for i := 0; i < pointCounts; i++ {
		points = append(points, point{
			Latitude:  center.Latitude + (rand.Float64()-0.5)*width,
			Longitude: center.Longitude + (rand.Float64()-0.5)*height,
		})
	}

	return convexHull(points)
}

func convexHull(p []point) []point {
	sort.Slice(p, func(i, j int) bool {
		if p[i].Latitude == p[j].Latitude {
			return p[i].Longitude < p[i].Longitude
		}
		return p[i].Latitude < p[j].Latitude
	})

	var h []point

	// Lower hull
	for _, pt := range p {
		for len(h) >= 2 && !ccw(h[len(h)-2], h[len(h)-1], pt) {
			h = h[:len(h)-1]
		}
		h = append(h, pt)
	}

	// Upper hull
	for i, t := len(p)-2, len(h)+1; i >= 0; i-- {
		pt := p[i]
		for len(h) >= t && !ccw(h[len(h)-2], h[len(h)-1], pt) {
			h = h[:len(h)-1]
		}
		h = append(h, pt)
	}

	return h[:len(h)-1]
}

// ccw returns true if the three points make a counter-clockwise turn
func ccw(a, b, c point) bool {
	return ((b.Latitude - a.Latitude) * (c.Longitude - a.Longitude)) > ((b.Longitude - a.Longitude) * (c.Latitude - a.Latitude))
}

func ToCoordinates(po []point) *client.Coordinates {
	r := make([]*client.Coordinate, 0, len(po)+1)

	for _, p := range po {
		r = append(r, &client.Coordinate{Latitude: p.Latitude, Longitude: p.Longitude})
	}

	// 始点と終点を一致させる
	r = append(r, r[0])

	return &client.Coordinates{Coordinates: r}
}

func contains(s []int64, e int64) bool {
	for _, v := range s {
		if e == v {
			return true
		}
	}
	return false
}

func getBoundingBox(points []point) [2]point {
	boundingBox := [2]point{
		{
			// TopLeftCorner
			Latitude: points[0].Latitude, Longitude: points[0].Longitude,
		},
		{
			// BottomRightCorner
			Latitude: points[0].Latitude, Longitude: points[0].Longitude,
		},
	}

	po := points[1:]

	for _, p := range po {
		if boundingBox[0].Latitude > p.Latitude {
			boundingBox[0].Latitude = p.Latitude
		}
		if boundingBox[0].Longitude > p.Longitude {
			boundingBox[0].Longitude = p.Longitude
		}

		if boundingBox[1].Latitude < p.Latitude {
			boundingBox[1].Latitude = p.Latitude
		}
		if boundingBox[1].Longitude < p.Longitude {
			boundingBox[1].Longitude = p.Longitude
		}
	}
	return boundingBox
}

func estateNazotteSearchScenario(ctx context.Context, c *client.Client) error {
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
	err = c.AccessEstateNazottePage(ctx)
	if err != nil {
		fails.Add(err)
		return failure.New(fails.ErrApplication)
	}
	if time.Since(t) > parameter.ThresholdTimeOfAbandonmentPage {
		return failure.New(fails.ErrTimeout)
	}

	// Nazotte Search
	// create nazotte data randomly
	convexHulled := createRandomConvexhull()
	polygon := ToCoordinates(convexHulled)
	boundingBox := getBoundingBox(convexHulled)

	t = time.Now()
	er, err := c.SearchEstatesNazotte(ctx, polygon)
	if err != nil {
		fails.Add(err)
		return failure.New(fails.ErrApplication)
	}

	if time.Since(t) > parameter.ThresholdTimeOfAbandonmentPage {
		return failure.New(fails.ErrTimeout)
	}

	if len(er.Estates) > parameter.MaxLengthOfNazotteResponse {
		err = failure.New(fails.ErrApplication, failure.Message("POST /api/estate/nazotte: レスポンスの内容が不正です"))
		fails.Add(err)
		return failure.New(fails.ErrApplication)
	}

	if err := checkEstatesInBoundingBox(er.Estates, boundingBox); err != nil {
		err = failure.Translate(err, fails.ErrApplication, failure.Message("GET /api/estate/nazotte: レスポンスの内容が不正です"))
		fails.Add(err)
		return failure.New(fails.ErrApplication)
	}

	if len(er.Estates) == 0 {
		return nil
	}

	randomPosition := rand.Intn(len(er.Estates))
	targetID := er.Estates[randomPosition].ID
	t = time.Now()
	e, err := c.AccessEstateDetailPage(ctx, targetID)
	if err != nil {
		fails.Add(err)
		return failure.New(fails.ErrApplication)
	}

	if time.Since(t) > parameter.ThresholdTimeOfAbandonmentPage {
		return failure.New(fails.ErrTimeout)
	}

	estate, err := asset.GetEstateFromID(e.ID)
	if err != nil || !e.Equal(estate) {
		err = failure.New(fails.ErrApplication, failure.Message("GET /api/estate/:id: レスポンスの内容が不正です"))
		fails.Add(err)
		return failure.New(fails.ErrApplication)
	}

	err = c.RequestEstateDocument(ctx, strconv.FormatInt(targetID, 10))
	if err != nil {
		fails.Add(err)
		return failure.New(fails.ErrApplication)
	}

	return nil
}
