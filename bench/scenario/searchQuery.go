package scenario

import (
	"math/rand"
	"net/url"
	"strconv"
	"strings"

	"github.com/isucon10-qualify/isucon10-qualify/bench/asset"
	"github.com/isucon10-qualify/isucon10-qualify/bench/parameter"
	"github.com/isucon10-qualify/isucon10-qualify/bench/score"
)

func randomTakeMany(slice []string, minLength, maxLength int) []string {
	s := make([]string, len(slice))
	copy(s, slice)
	rand.Shuffle(len(s), func(i, j int) { s[i], s[j] = s[j], s[i] })
	length := rand.Intn(maxLength-minLength) + minLength
	return s[:length]
}

func createRandomChairSearchQuery() (url.Values, error) {
	condition, err := asset.GetChairSearchCondition()
	if err != nil {
		return nil, err
	}

	level := score.GetLevel()
	paramNum := int(level/2 + 1)

	q := url.Values{}
	q.Set("perPage", strconv.Itoa(parameter.PerPageOfChairSearch))
	q.Set("page", "0")

	for i := 0; i < paramNum; i++ {
		r := rand.Intn(9)
		if level >= 2 {
			r += rand.Intn(1)
		}

		switch r {
		case 0, 1, 2:
			priceRangeID := condition.Price.Ranges[rand.Intn(len(condition.Price.Ranges))].ID
			q.Set("priceRangeId", strconv.FormatInt(priceRangeID, 10))

		case 3, 4:
			heightRangeID := condition.Height.Ranges[rand.Intn(len(condition.Height.Ranges))].ID
			q.Set("heightRangeId", strconv.FormatInt(heightRangeID, 10))

		case 5:
			widthRangeID := condition.Width.Ranges[rand.Intn(len(condition.Width.Ranges))].ID
			q.Set("widthRangeId", strconv.FormatInt(widthRangeID, 10))

		case 6:
			depthRangeID := condition.Depth.Ranges[rand.Intn(len(condition.Depth.Ranges))].ID
			q.Set("depthRangeId", strconv.FormatInt(depthRangeID, 10))

		case 7:
			kind := condition.Kind.List[rand.Intn(len(condition.Kind.List))]
			q.Set("kind", kind)

		case 8:
			color := condition.Color.List[rand.Intn(len(condition.Color.List))]
			q.Set("color", color)

		case 9:
			features := strings.Join(randomTakeMany(condition.Feature.List, 1, 3), ",")
			q.Set("features", features)
		}
	}

	return q, nil
}

func createRandomEstateSearchQuery() (url.Values, error) {
	condition, err := asset.GetEstateSearchCondition()
	if err != nil {
		return nil, err
	}

	level := score.GetLevel()
	paramNum := int(level/2 + 1)

	q := url.Values{}
	q.Set("perPage", strconv.Itoa(parameter.PerPageOfEstateSearch))
	q.Set("page", "0")

	for i := 0; i < paramNum; i++ {
		r := rand.Intn(5)
		if level >= 2 {
			r += rand.Intn(1)
		}

		switch r {
		case 0, 1, 2:
			rentRangeID := condition.Rent.Ranges[rand.Intn(len(condition.Rent.Ranges))].ID
			q.Set("rentRangeId", strconv.FormatInt(rentRangeID, 10))

		case 3:
			doorHeightRangeID := condition.DoorHeight.Ranges[rand.Intn(len(condition.DoorHeight.Ranges))].ID
			q.Set("doorHeightRangeId", strconv.FormatInt(doorHeightRangeID, 10))

		case 4:
			doorWidthRangeID := condition.DoorWidth.Ranges[rand.Intn(len(condition.DoorWidth.Ranges))].ID
			q.Set("doorWidthRangeId", strconv.FormatInt(doorWidthRangeID, 10))

		case 5:
			features := strings.Join(randomTakeMany(condition.Feature.List, 1, 3), ",")
			q.Set("features", features)
		}
	}

	return q, nil
}
