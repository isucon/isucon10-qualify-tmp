package scenario

import (
	"context"
	"path"
	"path/filepath"
	"strconv"
	"sync"

	"github.com/isucon10-qualify/isucon10-qualify/bench/asset"
	"github.com/isucon10-qualify/isucon10-qualify/bench/client"
	"github.com/isucon10-qualify/isucon10-qualify/bench/fails"
	"github.com/isucon10-qualify/isucon10-qualify/bench/parameter"
	"github.com/morikuni/failure"
	"golang.org/x/sync/errgroup"
)

// Verify Initialize後のアプリケーションサーバーに対して、副作用のない検証を実行する
// 早い段階でベンチマークをFailさせて早期リターンさせるのが目的
// ex) Search API を叩いて初期状態を確認する
func Verify(ctx context.Context, dataDir, fixtureDir string) {
	ctx, cancel := context.WithTimeout(ctx, parameter.VerifyTimeout)
	defer cancel()

	doneChan := make(chan bool)
	go func() {
		for {
			select {
			case <-ctx.Done():
				close(doneChan)
				return
			case <-fails.Fail():
			}
		}
	}()

	c := client.NewClientForVerify()

	verifyWithSnapshot(ctx, c, filepath.Join(dataDir, "result/verification_data"))
	if ctx.Err() != nil {
		err := failure.New(fails.ErrCritical, failure.Message("アプリケーション互換性チェックがタイムアウトしました"))
		fails.Add(err)
	}

	verifyWithScenario(ctx, c, fixtureDir, dataDir)
	if ctx.Err() != nil {
		err := failure.New(fails.ErrCritical, failure.Message("アプリケーション互換性チェックがタイムアウトしました"))
		fails.Add(err)
	}

	cancel()
	<-doneChan
	for {
		select {
		case <-fails.Fail():
		default:
			return
		}
	}
}

func verifyPostEstates(ctx context.Context, c *client.Client, estates []asset.Estate) error {
	id := strconv.FormatInt(estates[0].ID, 10)
	estate, err := c.GetEstateDetailFromID(ctx, id)
	if err == nil {
		return failure.Translate(err, fails.ErrApplication, failure.Message("未登録物件の詳細取得のレスポンスが不正です"))
	}

	err = c.PostEstates(ctx, estates)
	if err != nil {
		return failure.Translate(err, fails.ErrApplication, failure.Message("物件のCSV入稿に失敗しました"))
	}

	estate, err = c.GetEstateDetailFromID(ctx, id)
	if err != nil {
		return failure.Translate(err, fails.ErrApplication, failure.Message("登録済み物件の詳細取得に失敗しました"))
	}
	if estate == nil {
		return failure.New(fails.ErrApplication, failure.Message("登録済み物件の詳細取得に失敗しました"))
	}

	return nil
}

func verifyPostChairs(ctx context.Context, c *client.Client, chairs []asset.Chair) error {
	id := strconv.FormatInt(chairs[0].ID, 10)
	chair, err := c.GetChairDetailFromID(ctx, id)
	if err != nil {
		return failure.Translate(err, fails.ErrApplication, failure.Message("未登録イスの詳細取得のレスポンスが不正です"))
	}
	if chair != nil {
		return failure.New(fails.ErrApplication, failure.Message("未登録イスの詳細取得のレスポンスの内容が不正です"))
	}

	err = c.PostChairs(ctx, chairs)
	if err != nil {
		return failure.Translate(err, fails.ErrApplication, failure.Message("イスのCSV入稿に失敗しました"))
	}

	chair, err = c.GetChairDetailFromID(ctx, id)
	if err != nil {
		return failure.Translate(err, fails.ErrApplication, failure.Message("登録済みイスの詳細取得に失敗しました"))
	}
	if chair == nil {
		return failure.New(fails.ErrApplication, failure.Message("登録済みイスの詳細取得に失敗しました"))
	}

	return nil
}

func verifyChairStock(ctx context.Context, c *client.Client, id int64) error {
	strID := strconv.FormatInt(id, 10)

	chair, err := c.GetChairDetailFromID(ctx, strID)
	if err != nil {
		if ctxErr := ctx.Err(); ctxErr != nil {
			return ctxErr
		}
		return failure.Translate(err, fails.ErrApplication, failure.Message("イスの詳細取得に失敗しました"))
	}
	if chair == nil {
		return failure.New(fails.ErrApplication, failure.Message("在庫のあるはずのイスが売り切れになっています"))
	}

	err = c.BuyChair(ctx, strID)
	if err != nil {
		if ctxErr := ctx.Err(); ctxErr != nil {
			return ctxErr
		}
		return failure.Translate(err, fails.ErrApplication, failure.Message("イスの購入に失敗しました"))
	}

	chair, err = c.GetChairDetailFromID(ctx, strID)
	if err != nil {
		if ctxErr := ctx.Err(); ctxErr != nil {
			return ctxErr
		}
		return failure.Translate(err, fails.ErrApplication, failure.Message("売り切れたイスが存在します"))
	}

	if chair != nil {
		if ctxErr := ctx.Err(); ctxErr != nil {
			return ctxErr
		}
		return failure.New(fails.ErrApplication, failure.Message("売り切れたイスの詳細が表示されています"))
	}

	return nil
}

func verifyWithScenario(ctx context.Context, c *client.Client, fixtureDir, snapshotsParentsDirPath string) {
	var (
		estates []asset.Estate
		chairs  []asset.Chair
	)

	eg, childCtx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		var err error
		estates, err = loadEstatesFromJSON(childCtx, path.Join(snapshotsParentsDirPath, "result/verify_draft_estate.txt"))
		if err != nil {
			return err
		}

		return nil
	})

	eg.Go(func() error {
		var err error
		chairs, err = loadChairsFromJSON(childCtx, path.Join(snapshotsParentsDirPath, "result/verify_draft_chair.txt"))
		if err != nil {
			return err
		}

		return nil
	})

	if err := eg.Wait(); err != nil {
		fails.Add(failure.Translate(err, fails.ErrBenchmarker))
		return
	}

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := verifyPostEstates(ctx, c, estates)
		if err != nil {
			fails.Add(err)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		err := verifyPostChairs(ctx, c, chairs)
		if err != nil {
			fails.Add(err)
		}
	}()

	wg.Wait()

	wg.Add(1)
	go func() {
		defer wg.Done()
		err := verifyChairStock(ctx, c, chairs[0].ID)
		if err != nil {
			fails.Add(err)
		}
	}()

	wg.Wait()
}
