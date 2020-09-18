package scenario

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/isucon10-qualify/isucon10-qualify/bench/asset"
	"github.com/isucon10-qualify/isucon10-qualify/bench/client"
	"github.com/isucon10-qualify/isucon10-qualify/bench/fails"
	"github.com/morikuni/failure"
)

const (
	NumOfVerifyChairDetail                = 5
	NumOfVerifyChairSearchCondition       = 1
	NumOfVerifyChairSearch                = 5
	NumOfVerifyEstateDetail               = 5
	NumOfVerifyEstateSearchCondition      = 1
	NumOfVerifyEstateSearch               = 5
	NumOfVerifyLowPricedChair             = 1
	NumOfVerifyLowPricedEstate            = 1
	NumOfVerifyRecommendedEstateWithChair = 5
	NumOfVerifyEstateNazotte              = 5
)

var (
	ignoreChairUnexported  = cmpopts.IgnoreUnexported(asset.Chair{})
	ignoreEstateUnexported = cmpopts.IgnoreUnexported(asset.Estate{})
	ignoreEstateLatitude   = cmpopts.IgnoreFields(asset.Estate{}, "Latitude")
	ignoreEstateLongitude  = cmpopts.IgnoreFields(asset.Estate{}, "Longitude")
)

type Request struct {
	Method   string `json:"method"`
	Resource string `json:"resource"`
	Query    string `json:"query"`
	Body     string `json:"body"`
}

type Response struct {
	StatusCode int    `json:"statusCode"`
	Body       string `json:"body"`
}

type Snapshot struct {
	Request  Request  `json:"request"`
	Response Response `json:"response"`
}

func loadSnapshotFromFile(filePath string) (*Snapshot, error) {
	raw, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var snapshot *Snapshot
	err = json.Unmarshal(raw, &snapshot)
	if err != nil {
		return nil, err
	}

	return snapshot, nil
}

func verifyChairDetail(ctx context.Context, c *client.Client, filePath string) error {
	snapshot, err := loadSnapshotFromFile(filePath)
	if err != nil {
		return failure.Translate(err, fails.ErrBenchmarker, failure.Message("GET /api/chair/id: Snapshotの読み込みに失敗しました"))
	}

	idx := strings.LastIndex(snapshot.Request.Resource, "/")
	if idx == -1 || idx == len(snapshot.Request.Resource)-1 {
		return failure.Translate(err, fails.ErrBenchmarker, failure.Message("GET /api/chair/:id: 不正なSnapshotです"), failure.Messagef("snapshot: %s", filePath))
	}

	id := snapshot.Request.Resource[idx+1:]
	actual, err := c.GetChairDetailFromID(ctx, id)

	switch snapshot.Response.StatusCode {
	case http.StatusOK:
		if err != nil {
			return failure.Translate(err, fails.ErrApplication, failure.Message("GET /api/chair/:id: レスポンスが不正です"))
		}

		var expected *asset.Chair
		err = json.Unmarshal([]byte(snapshot.Response.Body), &expected)
		if err != nil {
			return failure.Translate(err, fails.ErrBenchmarker, failure.Message("GET /api/chair/:id: SnapshotのResponse BodyのUnmarshalでエラーが発生しました"), failure.Messagef("snapshot: %s", filePath))
		}

		if actual == nil {
			return failure.New(fails.ErrApplication, failure.Message("GET /api/chair/:id: レスポンスが不正です"), failure.Messagef("snapshot: %s", filePath))
		}

		if !cmp.Equal(*expected, *actual, ignoreChairUnexported) {
			return failure.New(fails.ErrApplication, failure.Message("GET /api/chair/:id: レスポンスが不正です"), failure.Messagef("snapshot: %s", filePath))
		}

	case http.StatusNotFound:
		if actual != nil {
			return failure.New(fails.ErrApplication, failure.Message("GET /api/chair/:id: レスポンスが不正です"))
		}
		if err != nil {
			return failure.Translate(err, fails.ErrApplication, failure.Message("GET /api/chair/:id: レスポンスが不正です"))
		}

	default:
		if err == nil {
			return failure.New(fails.ErrApplication, failure.Message("GET /api/chair/:id: レスポンスが不正です"))
		}
	}

	return nil
}

func verifyChairSearchCondition(ctx context.Context, c *client.Client, filePath string) error {
	snapshot, err := loadSnapshotFromFile(filePath)
	if err != nil {
		return failure.Translate(err, fails.ErrBenchmarker, failure.Message("GET /api/chair/search/condition: Snapshotの読み込みに失敗しました"))
	}

	actual, err := c.GetChairSearchCondition(ctx)

	switch snapshot.Response.StatusCode {
	case http.StatusOK:
		if err != nil {
			return failure.Translate(err, fails.ErrApplication, failure.Message("GET /api/chair/search/condition: レスポンスが不正です"))
		}

		var expected *asset.ChairSearchCondition
		err = json.Unmarshal([]byte(snapshot.Response.Body), &expected)
		if err != nil {
			return failure.Translate(err, fails.ErrBenchmarker, failure.Message("GET /api/chair/search/condition: SnapshotのResponse BodyのUnmarshalでエラーが発生しました"), failure.Messagef("snapshot: %s", filePath))
		}

		if !cmp.Equal(*expected, *actual, ignoreChairUnexported) {
			return failure.New(fails.ErrApplication, failure.Message("GET /api/chair/search/condition: レスポンスが不正です"), failure.Messagef("snapshot: %s", filePath))
		}

	default:
		if err == nil {
			return failure.New(fails.ErrApplication, failure.Message("GET /api/chair/search/condition: レスポンスが不正です"))
		}
	}

	return nil
}

func verifyChairSearch(ctx context.Context, c *client.Client, filePath string) error {
	snapshot, err := loadSnapshotFromFile(filePath)
	if err != nil {
		return failure.Translate(err, fails.ErrBenchmarker, failure.Message("GET /api/chair/search: Snapshotの読み込みに失敗しました"))
	}

	q, err := url.ParseQuery(snapshot.Request.Query)
	if err != nil {
		return failure.Translate(err, fails.ErrBenchmarker, failure.Message("GET /api/chair/search: Request QueryのUnmarshalでエラーが発生しました"))
	}

	actual, err := c.SearchChairsWithQuery(ctx, q)

	switch snapshot.Response.StatusCode {
	case http.StatusOK:
		if err != nil {
			return failure.Translate(err, fails.ErrApplication, failure.Message("GET /api/chair/search: レスポンスが不正です"))
		}

		var expected *client.ChairsResponse
		err = json.Unmarshal([]byte(snapshot.Response.Body), &expected)
		if err != nil {
			return failure.Translate(err, fails.ErrBenchmarker, failure.Message("GET /api/chair/search: SnapshotのResponse BodyのUnmarshalでエラーが発生しました"), failure.Messagef("snapshot: %s", filePath))
		}

		if !cmp.Equal(*expected, *actual, ignoreChairUnexported) {
			return failure.New(fails.ErrApplication, failure.Message("GET /api/chair/search: レスポンスが不正です"), failure.Messagef("snapshot: %s", filePath))
		}

	default:
		if err == nil {
			return failure.New(fails.ErrApplication, failure.Message("GET /api/chair/search: レスポンスが不正です"))
		}
	}

	return nil
}

func verifyEstateDetail(ctx context.Context, c *client.Client, filePath string) error {
	snapshot, err := loadSnapshotFromFile(filePath)
	if err != nil {
		return failure.Translate(err, fails.ErrBenchmarker, failure.Message("GET /api/estate/id: Snapshotの読み込みに失敗しました"))
	}

	idx := strings.LastIndex(snapshot.Request.Resource, "/")
	if idx == -1 || idx == len(snapshot.Request.Resource)-1 {
		return failure.Translate(err, fails.ErrBenchmarker, failure.Message("GET /api/estate/:id: 不正なSnapshotです"), failure.Messagef("snapshot: %s", filePath))
	}

	id := snapshot.Request.Resource[idx+1:]
	actual, err := c.GetEstateDetailFromID(ctx, id)

	switch snapshot.Response.StatusCode {
	case http.StatusOK:
		if err != nil {
			return failure.Translate(err, fails.ErrApplication, failure.Message("GET /api/estate/:id: レスポンスが不正です"))
		}

		var expected *asset.Estate
		err = json.Unmarshal([]byte(snapshot.Response.Body), &expected)
		if err != nil {
			return failure.Translate(err, fails.ErrBenchmarker, failure.Message("GET /api/estate/:id: SnapshotのResponse BodyのUnmarshalでエラーが発生しました"), failure.Messagef("snapshot: %s", filePath))
		}

		if !cmp.Equal(*expected, *actual, ignoreEstateUnexported, ignoreEstateLatitude, ignoreEstateLongitude) {
			return failure.New(fails.ErrApplication, failure.Message("GET /api/estate/:id: レスポンスが不正です"), failure.Messagef("snapshot: %s", filePath))
		}

	default:
		if err == nil {
			return failure.New(fails.ErrApplication, failure.Message("GET /api/estate/:id: レスポンスが不正です"))
		}
	}

	return nil
}

func verifyEstateSearchCondition(ctx context.Context, c *client.Client, filePath string) error {
	snapshot, err := loadSnapshotFromFile(filePath)
	if err != nil {
		return failure.Translate(err, fails.ErrBenchmarker, failure.Message("GET /api/estate/search/condition: Snapshotの読み込みに失敗しました"))
	}

	actual, err := c.GetEstateSearchCondition(ctx)

	switch snapshot.Response.StatusCode {
	case http.StatusOK:
		if err != nil {
			return failure.Translate(err, fails.ErrApplication, failure.Message("GET /api/estate/search/condition: レスポンスが不正です"))
		}

		var expected *asset.EstateSearchCondition
		err = json.Unmarshal([]byte(snapshot.Response.Body), &expected)
		if err != nil {
			return failure.Translate(err, fails.ErrBenchmarker, failure.Message("GET /api/estate/search/condition: SnapshotのResponse BodyのUnmarshalでエラーが発生しました"), failure.Messagef("snapshot: %s", filePath))
		}

		if !cmp.Equal(*expected, *actual) {
			return failure.New(fails.ErrApplication, failure.Message("GET /api/estate/search/condition: レスポンスが不正です"), failure.Messagef("snapshot: %s", filePath))
		}

	default:
		if err == nil {
			return failure.New(fails.ErrApplication, failure.Message("GET /api/estate/search/condition: レスポンスが不正です"))
		}
	}

	return nil
}

func verifyEstateSearch(ctx context.Context, c *client.Client, filePath string) error {
	snapshot, err := loadSnapshotFromFile(filePath)
	if err != nil {
		return failure.Translate(err, fails.ErrBenchmarker, failure.Message("GET /api/estate/search: Snapshotの読み込みに失敗しました"))
	}

	q, err := url.ParseQuery(snapshot.Request.Query)
	if err != nil {
		return failure.Translate(err, fails.ErrBenchmarker, failure.Message("GET /api/estate/search: Request QueryのUnmarshalでエラーが発生しました"))
	}

	actual, err := c.SearchEstatesWithQuery(ctx, q)

	switch snapshot.Response.StatusCode {
	case http.StatusOK:
		if err != nil {
			return failure.Translate(err, fails.ErrApplication, failure.Message("GET /api/estate/search: レスポンスが不正です"))
		}

		var expected *client.EstatesResponse
		err = json.Unmarshal([]byte(snapshot.Response.Body), &expected)
		if err != nil {
			return failure.Translate(err, fails.ErrBenchmarker, failure.Message("GET /api/estate/search: SnapshotのResponse BodyのUnmarshalでエラーが発生しました"), failure.Messagef("snapshot: %s", filePath))
		}

		if !cmp.Equal(*expected, *actual, ignoreEstateUnexported, ignoreEstateLatitude, ignoreEstateLongitude) {
			return failure.New(fails.ErrApplication, failure.Message("GET /api/estate/search: レスポンスが不正です"), failure.Messagef("snapshot: %s", filePath))
		}

	default:
		if err == nil {
			return failure.New(fails.ErrApplication, failure.Message("GET /api/estate/search: レスポンスが不正です"))
		}
	}

	return nil
}

func verifyLowPricedChair(ctx context.Context, c *client.Client, filePath string) error {
	snapshot, err := loadSnapshotFromFile(filePath)
	if err != nil {
		return failure.Translate(err, fails.ErrBenchmarker, failure.Message("GET /api/chair/low_priced: Snapshotの読み込みに失敗しました"))
	}

	actual, err := c.GetLowPricedChair(ctx)

	switch snapshot.Response.StatusCode {
	case http.StatusOK:
		if err != nil {
			return failure.Translate(err, fails.ErrApplication, failure.Message("GET /api/chair/low_priced: レスポンスが不正です"))
		}

		var expected *client.ChairsResponse
		err = json.Unmarshal([]byte(snapshot.Response.Body), &expected)
		if err != nil {
			return failure.Translate(err, fails.ErrBenchmarker, failure.Message("GET /api/chair/low_priced: SnapshotのResponse BodyのUnmarshalでエラーが発生しました"), failure.Messagef("snapshot: %s", filePath))
		}

		if !cmp.Equal(*expected, *actual, ignoreChairUnexported) {
			return failure.New(fails.ErrApplication, failure.Message("GET /api/chair/low_priced: レスポンスが不正です"), failure.Messagef("snapshot: %s", filePath))
		}

	default:
		if err == nil {
			return failure.New(fails.ErrApplication, failure.Message("GET /api/chair/low_priced: レスポンスが不正です"))
		}
	}

	return nil
}

func verifyLowPricedEstate(ctx context.Context, c *client.Client, filePath string) error {
	snapshot, err := loadSnapshotFromFile(filePath)
	if err != nil {
		return failure.Translate(err, fails.ErrBenchmarker, failure.Message("GET /api/estate/low_priced: Snapshotの読み込みに失敗しました"))
	}

	actual, err := c.GetLowPricedEstate(ctx)

	switch snapshot.Response.StatusCode {
	case http.StatusOK:
		if err != nil {
			return failure.Translate(err, fails.ErrApplication, failure.Message("GET /api/estate/low_priced: レスポンスが不正です"))
		}

		var expected *client.EstatesResponse
		err = json.Unmarshal([]byte(snapshot.Response.Body), &expected)
		if err != nil {
			return failure.Translate(err, fails.ErrBenchmarker, failure.Message("GET /api/estate/low_priced: SnapshotのResponse BodyのUnmarshalでエラーが発生しました"), failure.Messagef("snapshot: %s", filePath))
		}

		if !cmp.Equal(*expected, *actual, ignoreEstateUnexported, ignoreEstateLatitude, ignoreEstateLongitude) {
			return failure.New(fails.ErrApplication, failure.Message("GET /api/estate/low_priced: レスポンスが不正です"), failure.Messagef("snapshot: %s", filePath))
		}

	default:
		if err == nil {
			return failure.New(fails.ErrApplication, failure.Message("GET /api/estate/low_priced: レスポンスが不正です"))
		}
	}

	return nil
}

func verifyRecommendedEstateWithChair(ctx context.Context, c *client.Client, filePath string) error {
	snapshot, err := loadSnapshotFromFile(filePath)
	if err != nil {
		return failure.Translate(err, fails.ErrBenchmarker, failure.Message("GET /api/recommended_estate/:id: Snapshotの読み込みに失敗しました"))
	}

	idx := strings.LastIndex(snapshot.Request.Resource, "/")
	if idx == -1 || idx == len(snapshot.Request.Resource)-1 {
		return failure.Translate(err, fails.ErrBenchmarker, failure.Message("GET /api/recommended_estate/:id: 不正なSnapshotです"), failure.Messagef("snapshot: %s", filePath))
	}
	id, err := strconv.ParseInt(snapshot.Request.Resource[idx+1:], 10, 64)
	if err != nil {
		return failure.Translate(err, fails.ErrBenchmarker, failure.Message("GET /api/recommended_estate/:id: 不正なSnapshotです"), failure.Messagef("snapshot: %s", filePath))
	}

	actual, err := c.GetRecommendedEstatesFromChair(ctx, id)

	switch snapshot.Response.StatusCode {
	case http.StatusOK:
		if err != nil {
			return failure.Translate(err, fails.ErrApplication, failure.Message("GET /api/recommended_estate/:id: レスポンスが不正です"))
		}

		var expected *client.EstatesResponse
		err = json.Unmarshal([]byte(snapshot.Response.Body), &expected)
		if err != nil {
			return failure.Translate(err, fails.ErrBenchmarker, failure.Message("GET /api/recommended_estate/:id: SnapshotのResponse BodyのUnmarshalでエラーが発生しました"), failure.Messagef("snapshot: %s", filePath))
		}
		if !cmp.Equal(*expected, *actual, ignoreEstateUnexported, ignoreEstateLatitude, ignoreEstateLongitude) {
			return failure.New(fails.ErrApplication, failure.Message("GET /api/recommended_estate/:id: レスポンスが不正です"), failure.Messagef("snapshot: %s", filePath))
		}

	default:
		if err == nil {
			return failure.New(fails.ErrApplication, failure.Message("GET /api/recommended_estate/:id: レスポンスが不正です"))
		}
	}

	return nil
}

func verifyEstateNazotte(ctx context.Context, c *client.Client, filePath string) error {
	snapshot, err := loadSnapshotFromFile(filePath)
	if err != nil {
		return failure.Translate(err, fails.ErrBenchmarker, failure.Message("POST /api/estate/nazotte: Snapshotの読み込みに失敗しました"))
	}

	var coordinates *client.Coordinates
	err = json.Unmarshal([]byte(snapshot.Request.Body), &coordinates)
	if err != nil {
		return failure.Translate(err, fails.ErrBenchmarker, failure.Message("POST /api/estate/nazotte: Request BodyのUnmarshalでエラーが発生しました"))
	}

	actual, err := c.SearchEstatesNazotte(ctx, coordinates)

	switch snapshot.Response.StatusCode {
	case http.StatusOK:
		if err != nil {
			return failure.Translate(err, fails.ErrApplication, failure.Message("POST /api/estate/nazotte: レスポンスが不正です"))
		}

		var expected *client.EstatesResponse
		err = json.Unmarshal([]byte(snapshot.Response.Body), &expected)
		if err != nil {
			return failure.Translate(err, fails.ErrBenchmarker, failure.Message("POST /api/estate/nazotte: SnapshotのResponse BodyのUnmarshalでエラーが発生しました"), failure.Messagef("snapshot: %s", filePath))
		}

		if !cmp.Equal(*expected, *actual, ignoreEstateUnexported, ignoreEstateLatitude, ignoreEstateLongitude) {
			return failure.New(fails.ErrApplication, failure.Message("POST /api/estate/nazotte: レスポンスが不正です"), failure.Messagef("snapshot: %s", filePath))
		}

	default:
		if err == nil {
			return failure.New(fails.ErrApplication, failure.Message("POST /api/estate/nazotte: レスポンスが不正です"))
		}
	}

	return nil
}

func verifyWithSnapshot(ctx context.Context, c *client.Client, snapshotsParentsDirPath string) {
	wg := sync.WaitGroup{}

	snapshotsDirPath := filepath.Join(snapshotsParentsDirPath, "chair_detail")
	snapshots, err := ioutil.ReadDir(snapshotsDirPath)
	if err != nil {
		err := failure.Translate(err, fails.ErrBenchmarker, failure.Message("GET /api/chair/:id: Snapshotディレクトリがありません"))
		fails.Add(err)
	} else {
		for i := 0; i < NumOfVerifyChairDetail; i++ {
			wg.Add(1)
			r := rand.Intn(len(snapshots))
			go func(filePath string) {
				err := verifyChairDetail(ctx, c, filePath)
				if err != nil {
					fails.Add(err)
				}
				wg.Done()
			}(path.Join(snapshotsDirPath, snapshots[r].Name()))
		}
	}

	snapshotsDirPath = filepath.Join(snapshotsParentsDirPath, "chair_search_condition")
	snapshots, err = ioutil.ReadDir(snapshotsDirPath)
	if err != nil {
		err := failure.Translate(err, fails.ErrBenchmarker, failure.Message("GET /api/chair/search/condition: Snapshotディレクトリがありません"))
		fails.Add(err)
	} else {
		for i := 0; i < NumOfVerifyChairSearchCondition; i++ {
			wg.Add(1)
			r := rand.Intn(len(snapshots))
			go func(filePath string) {
				err := verifyChairSearchCondition(ctx, c, filePath)
				if err != nil {
					fails.Add(err)
				}
				wg.Done()
			}(path.Join(snapshotsDirPath, snapshots[r].Name()))
		}
	}

	snapshotsDirPath = filepath.Join(snapshotsParentsDirPath, "chair_search")
	snapshots, err = ioutil.ReadDir(snapshotsDirPath)
	if err != nil {
		err := failure.Translate(err, fails.ErrBenchmarker, failure.Message("GET /api/chair/search: Snapshotディレクトリがありません"))
		fails.Add(err)
	} else {
		for i := 0; i < NumOfVerifyChairSearch; i++ {
			wg.Add(1)
			r := rand.Intn(len(snapshots))
			go func(filePath string) {
				err := verifyChairSearch(ctx, c, filePath)
				if err != nil {
					fails.Add(err)
				}
				wg.Done()
			}(path.Join(snapshotsDirPath, snapshots[r].Name()))
		}
	}

	snapshotsDirPath = filepath.Join(snapshotsParentsDirPath, "estate_detail")
	snapshots, err = ioutil.ReadDir(snapshotsDirPath)
	if err != nil {
		err := failure.Translate(err, fails.ErrBenchmarker, failure.Message("GET /api/estate/:id: Snapshotディレクトリがありません"))
		fails.Add(err)
	} else {
		for i := 0; i < NumOfVerifyEstateDetail; i++ {
			wg.Add(1)
			r := rand.Intn(len(snapshots))
			go func(filePath string) {
				err := verifyEstateDetail(ctx, c, filePath)
				if err != nil {
					fails.Add(err)
				}
				wg.Done()
			}(path.Join(snapshotsDirPath, snapshots[r].Name()))
		}
	}

	snapshotsDirPath = filepath.Join(snapshotsParentsDirPath, "estate_search_condition")
	snapshots, err = ioutil.ReadDir(snapshotsDirPath)
	if err != nil {
		err := failure.Translate(err, fails.ErrBenchmarker, failure.Message("GET /api/estate/search/condition: Snapshotディレクトリがありません"))
		fails.Add(err)
	} else {
		for i := 0; i < NumOfVerifyEstateSearchCondition; i++ {
			wg.Add(1)
			r := rand.Intn(len(snapshots))
			go func(filePath string) {
				err := verifyEstateSearchCondition(ctx, c, filePath)
				if err != nil {
					fails.Add(err)
				}
				wg.Done()
			}(path.Join(snapshotsDirPath, snapshots[r].Name()))
		}
	}

	snapshotsDirPath = filepath.Join(snapshotsParentsDirPath, "estate_search")
	snapshots, err = ioutil.ReadDir(snapshotsDirPath)
	if err != nil {
		err := failure.Translate(err, fails.ErrBenchmarker, failure.Message("GET /api/estate/search: Snapshotディレクトリがありません"))
		fails.Add(err)
	} else {
		for i := 0; i < NumOfVerifyEstateSearch; i++ {
			wg.Add(1)
			r := rand.Intn(len(snapshots))
			go func(filePath string) {
				err := verifyEstateSearch(ctx, c, filePath)
				if err != nil {
					fails.Add(err)
				}
				wg.Done()
			}(path.Join(snapshotsDirPath, snapshots[r].Name()))
		}
	}

	snapshotsDirPath = filepath.Join(snapshotsParentsDirPath, "chair_low_priced")
	snapshots, err = ioutil.ReadDir(snapshotsDirPath)
	if err != nil {
		err := failure.Translate(err, fails.ErrBenchmarker, failure.Message("GET /api/chair/low_priced: Snapshotディレクトリがありません"))
		fails.Add(err)
	} else {
		for i := 0; i < NumOfVerifyLowPricedChair; i++ {
			wg.Add(1)
			r := rand.Intn(len(snapshots))
			go func(filePath string) {
				err := verifyLowPricedChair(ctx, c, filePath)
				if err != nil {
					fails.Add(err)
				}
				wg.Done()
			}(path.Join(snapshotsDirPath, snapshots[r].Name()))
		}
	}

	snapshotsDirPath = filepath.Join(snapshotsParentsDirPath, "estate_low_priced")
	snapshots, err = ioutil.ReadDir(snapshotsDirPath)
	if err != nil {
		err := failure.Translate(err, fails.ErrBenchmarker, failure.Message("GET /api/estate/low_priced: Snapshotディレクトリがありません"))
		fails.Add(err)
	} else {
		for i := 0; i < NumOfVerifyLowPricedEstate; i++ {
			wg.Add(1)
			r := rand.Intn(len(snapshots))
			go func(filePath string) {
				err := verifyLowPricedEstate(ctx, c, filePath)
				if err != nil {
					fails.Add(err)
				}
				wg.Done()
			}(path.Join(snapshotsDirPath, snapshots[r].Name()))
		}
	}

	snapshotsDirPath = filepath.Join(snapshotsParentsDirPath, "recommended_estate_with_chair")
	snapshots, err = ioutil.ReadDir(snapshotsDirPath)
	if err != nil {
		err := failure.Translate(err, fails.ErrBenchmarker, failure.Message("GET /api/recommended_estate/:id: Snapshotディレクトリがありません"))
		fails.Add(err)
	} else {
		for i := 0; i < NumOfVerifyRecommendedEstateWithChair; i++ {
			wg.Add(1)
			r := rand.Intn(len(snapshots))
			go func(filePath string) {
				err := verifyRecommendedEstateWithChair(ctx, c, filePath)
				if err != nil {
					fails.Add(err)
				}
				wg.Done()
			}(path.Join(snapshotsDirPath, snapshots[r].Name()))
		}
	}

	snapshotsDirPath = filepath.Join(snapshotsParentsDirPath, "estate_nazotte")
	snapshots, err = ioutil.ReadDir(snapshotsDirPath)
	if err != nil {
		err := failure.Translate(err, fails.ErrBenchmarker, failure.Message("POST /api/estate/nazotte: Snapshotディレクトリがありません"))
		fails.Add(err)
	} else {
		for i := 0; i < NumOfVerifyEstateNazotte; i++ {
			wg.Add(1)
			r := rand.Intn(len(snapshots))
			go func(filePath string) {
				err := verifyEstateNazotte(ctx, c, filePath)
				if err != nil {
					fails.Add(err)
				}
				wg.Done()
			}(path.Join(snapshotsDirPath, snapshots[r].Name()))
		}
	}

	wg.Wait()
}
