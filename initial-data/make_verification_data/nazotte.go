package main

import (
	"encoding/json"
	"math/rand"
	"sort"
)

var FamousPlaces []Coordinate = []Coordinate{
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

func convexHull(p []Coordinate) []Coordinate {
	sort.Slice(p, func(i, j int) bool {
		if p[i].Latitude == p[j].Latitude {
			return p[i].Longitude < p[i].Longitude
		}
		return p[i].Latitude < p[j].Latitude
	})

	var h []Coordinate

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

// ccw returns true if the three Coordinates make a counter-clockwise turn
func ccw(a, b, c Coordinate) bool {
	return ((b.Latitude - a.Latitude) * (c.Longitude - a.Longitude)) > ((b.Longitude - a.Longitude) * (c.Latitude - a.Latitude))
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

func createRandomConvexhull() string {
	famousPlace := FamousPlaces[rand.Intn(len(FamousPlaces))]

	width := rand.Float64()*(rangeMaxWidth-rangeMinWidth) + rangeMinWidth
	height := rand.Float64()*(rangeMaxHeight-rangeMinHeight) + rangeMinHeight
	center := Coordinate{
		Latitude:  famousPlace.Latitude + (rand.Float64()-0.5)*rangeDiffLatitude,
		Longitude: famousPlace.Longitude + (rand.Float64()-0.5)*rangeDiffLongitude,
	}

	pointCounts := rand.Intn(numOfMaxPoints-numOfMinPoints) + numOfMinPoints

	coordinates := []Coordinate{}

	for i := 0; i < pointCounts; i++ {
		coordinates = append(coordinates, Coordinate{
			Latitude:  center.Latitude + (rand.Float64()-0.5)*width,
			Longitude: center.Longitude + (rand.Float64()-0.5)*height,
		})
	}

	convexhulled := convexHull(coordinates)
	convexhulled = append(convexhulled, convexhulled[0])

	body, err := json.Marshal(NazotteRequestBody{
		Coordinates: convexhulled,
	})
	if err != nil {
		panic(err)
	}

	return string(body)
}
