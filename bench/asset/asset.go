package asset

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"os"
	"path/filepath"
	"sync"

	"github.com/isucon10-qualify/isucon10-qualify/bench/fails"
	"github.com/morikuni/failure"
	"golang.org/x/sync/errgroup"
)

var (
	chairMap  map[int64]*Chair
	chairMu   sync.RWMutex
	estateMap map[int64]*Estate
	estateMu  sync.RWMutex
)

// メモリ上にデータを展開する
// このデータを使用してAPIからのレスポンスを確認する
func Initialize(ctx context.Context, dataDir, fixtureDir string) {
	eg, childCtx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		f, err := os.Open(filepath.Join(dataDir, "result/chair_json.txt"))
		if err != nil {
			return err
		}
		defer f.Close()

		chairMap = map[int64]*Chair{}
		decoder := json.NewDecoder(f)
		for {
			if err := childCtx.Err(); err != nil {
				return err
			}

			var chair Chair
			if err := decoder.Decode(&chair); err != nil {
				if err == io.EOF {
					break
				}
				return err
			}
			StoreChair(chair)
		}
		return nil
	})

	eg.Go(func() error {
		f, err := os.Open(filepath.Join(dataDir, "result/estate_json.txt"))
		if err != nil {
			return err
		}
		defer f.Close()

		estateMap = map[int64]*Estate{}
		decoder := json.NewDecoder(f)
		for {
			if err := childCtx.Err(); err != nil {
				return err
			}

			var estate Estate
			if err := decoder.Decode(&estate); err != nil {
				if err == io.EOF {
					break
				}
				return err
			}
			StoreEstate(estate)
		}

		return nil
	})

	eg.Go(func() error {
		err := loadChairSearchCondition(fixtureDir)
		if err != nil {
			return err
		}
		return nil
	})

	eg.Go(func() error {
		err := loadEstateSearchCondition(fixtureDir)
		if err != nil {
			return err
		}
		return nil
	})

	eg.Go(func() error {
		err := loadChairDraftFiles(dataDir)
		if err != nil {
			return err
		}
		return nil
	})

	eg.Go(func() error {
		err := loadEstateDraftFiles(dataDir)
		if err != nil {
			return err
		}
		return nil
	})

	if err := eg.Wait(); err != nil {
		err = failure.Translate(err, fails.ErrBenchmarker, failure.Message("assetの初期化に失敗しました"))
		fails.Add(err)
	}
}

func GetChairFromID(id int64) (*Chair, error) {
	chairMu.RLock()
	defer chairMu.RUnlock()
	c, ok := chairMap[id]
	if !ok {
		return nil, errors.New("requested chair not found")
	}
	return c, nil
}

func StoreChair(chair Chair) {
	chairMu.Lock()
	defer chairMu.Unlock()
	chairMap[chair.ID] = &chair
}

func DecrementChairStock(id int64) {
	chairMu.RLock()
	defer chairMu.RUnlock()
	c, ok := chairMap[id]
	if ok {
		c.DecrementStock()
	}
}

func GetEstateFromID(id int64) (*Estate, error) {
	estateMu.RLock()
	defer estateMu.RUnlock()
	e, ok := estateMap[id]
	if !ok {
		return nil, errors.New("requested estate not found")
	}
	return e, nil
}

func StoreEstate(estate Estate) {
	estateMu.Lock()
	defer estateMu.Unlock()
	estateMap[estate.ID] = &estate
}
