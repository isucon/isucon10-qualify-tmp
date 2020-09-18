package main

import (
	"math/rand"
	"net/url"
	"strconv"
	"strings"
)

func createRandomEstateSearchQuery(condition EstateSearchCondition) url.Values {
	q := url.Values{}
	if (rand.Intn(100) % 10) == 0 {
		rentRangeID := condition.Rent.Ranges[rand.Intn(len(condition.Rent.Ranges))].ID
		q.Set("rentRangeId", strconv.FormatInt(rentRangeID, 10))
	}
	if (rand.Intn(100) % 10) == 0 {
		doorHeightRangeID := condition.DoorHeight.Ranges[rand.Intn(len(condition.DoorHeight.Ranges))].ID
		q.Set("doorHeightRangeId", strconv.FormatInt(doorHeightRangeID, 10))
	}
	if (rand.Intn(100) % 10) == 0 {
		doorWidthRangeID := condition.DoorWidth.Ranges[rand.Intn(len(condition.DoorWidth.Ranges))].ID
		q.Set("doorWidthRangeId", strconv.FormatInt(doorWidthRangeID, 10))
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
