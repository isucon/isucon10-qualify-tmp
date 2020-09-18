package asset

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"

	"github.com/morikuni/failure"

	"github.com/isucon10-qualify/isucon10-qualify/bench/fails"
)

var (
	chairSearchCondition  *ChairSearchCondition
	estateSearchCondition *EstateSearchCondition
)

type Range struct {
	ID  int64 `json:"id"`
	Min int64 `json:"min"`
	Max int64 `json:"max"`
}

type RangeCondition struct {
	Prefix string   `json:"prefix"`
	Suffix string   `json:"suffix"`
	Ranges []*Range `json:"ranges"`
}

type ListCondition struct {
	List []string `json:"list"`
}

type EstateSearchCondition struct {
	DoorWidth  RangeCondition `json:"doorWidth"`
	DoorHeight RangeCondition `json:"doorHeight"`
	Rent       RangeCondition `json:"rent"`
	Feature    ListCondition  `json:"feature"`
}

type ChairSearchCondition struct {
	Width   RangeCondition `json:"width"`
	Height  RangeCondition `json:"height"`
	Depth   RangeCondition `json:"depth"`
	Price   RangeCondition `json:"price"`
	Color   ListCondition  `json:"color"`
	Feature ListCondition  `json:"feature"`
	Kind    ListCondition  `json:"kind"`
}

func loadChairSearchCondition(fixtureDir string) error {
	jsonText, err := ioutil.ReadFile(filepath.Join(fixtureDir, "chair_condition.json"))
	if err != nil {
		return err
	}

	json.Unmarshal(jsonText, &chairSearchCondition)
	return nil
}

func loadEstateSearchCondition(fixtureDir string) error {
	jsonText, err := ioutil.ReadFile(filepath.Join(fixtureDir, "estate_condition.json"))
	if err != nil {
		return err
	}

	json.Unmarshal(jsonText, &estateSearchCondition)
	return nil
}

func GetChairSearchCondition() (*ChairSearchCondition, error) {
	if chairSearchCondition == nil {
		return nil, failure.New(fails.ErrBenchmarker, failure.Message("イスの検索条件が読み込まれていません"))
	}
	return chairSearchCondition, nil
}

func GetEstateSearchCondition() (*EstateSearchCondition, error) {
	if estateSearchCondition == nil {
		return nil, failure.New(fails.ErrBenchmarker, failure.Message("物件の検索条件が読み込まれていません"))
	}
	return estateSearchCondition, nil
}
