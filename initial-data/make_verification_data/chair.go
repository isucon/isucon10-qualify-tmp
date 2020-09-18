package main

import (
	"math/rand"
	"net/url"
	"strconv"
	"strings"
)

func createRandomChairSearchQuery(condition ChairSearchCondition) url.Values {
	q := url.Values{}
	priceRangeID := condition.Price.Ranges[rand.Intn(len(condition.Price.Ranges))].ID
	if (rand.Intn(100) % 10) == 0 {
		q.Set("priceRangeId", strconv.FormatInt(priceRangeID, 10))
	}
	if (rand.Intn(100) % 10) == 0 {
		heightRangeID := condition.Height.Ranges[rand.Intn(len(condition.Height.Ranges))].ID
		q.Set("heightRangeId", strconv.FormatInt(heightRangeID, 10))
	}
	if (rand.Intn(100) % 10) == 0 {
		widthRangeID := condition.Width.Ranges[rand.Intn(len(condition.Width.Ranges))].ID
		q.Set("widthRangeId", strconv.FormatInt(widthRangeID, 10))
	}
	if (rand.Intn(100) % 10) == 0 {
		depthRangeID := condition.Depth.Ranges[rand.Intn(len(condition.Depth.Ranges))].ID
		q.Set("depthRangeId", strconv.FormatInt(depthRangeID, 10))
	}

	if (rand.Intn(100) % 10) == 0 {
		q.Set("kind", condition.Kind.List[rand.Intn(len(condition.Kind.List))])
	}
	if (rand.Intn(100) % 10) == 0 {
		q.Set("color", condition.Color.List[rand.Intn(len(condition.Color.List))])
	}
	// condition.Featureの最後の1つはScenario形式のVerify用で該当件数が少ないため、Snapshot形式のVerifyでは使用しない
	features := make([]string, len(condition.Feature.List)-1)
	copy(features, condition.Feature.List[:len(condition.Feature.List)-1])
	rand.Shuffle(len(features), func(i, j int) { features[i], features[j] = features[j], features[i] })
	featureLength := rand.Intn(len(features)) + 1
	if featureLength > 3 {
		featureLength = rand.Intn(3) + 1
	}
	q.Set("features", strings.Join(features[:featureLength], ","))

	q.Set("perPage", strconv.Itoa(rand.Intn(30)+20))
	q.Set("page", "0")

	return q
}
