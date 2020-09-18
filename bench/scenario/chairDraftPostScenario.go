package scenario

import (
	"context"
	"encoding/json"
	"io"
	"os"
	"strconv"

	"github.com/isucon10-qualify/isucon10-qualify/bench/asset"
	"github.com/isucon10-qualify/isucon10-qualify/bench/client"
	"github.com/isucon10-qualify/isucon10-qualify/bench/fails"
	"github.com/morikuni/failure"
)

func loadChairsFromJSON(ctx context.Context, filePath string) ([]asset.Chair, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	chairs := []asset.Chair{}
	decoder := json.NewDecoder(f)
	for {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		var chair asset.Chair
		if err := decoder.Decode(&chair); err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		chairs = append(chairs, chair)
	}

	return chairs, nil
}

func chairDraftPostScenario(ctx context.Context, c *client.Client, filePath string) {
	chairs, err := loadChairsFromJSON(ctx, filePath)
	if err != nil {
		fails.Add(failure.Translate(err, fails.ErrCritical))
		return
	}

	id := strconv.FormatInt(chairs[0].ID, 10)
	chair, err := c.GetChairDetailFromID(ctx, id)
	if err != nil {
		fails.Add(failure.Translate(err, fails.ErrCritical))
		return
	}
	if chair != nil {
		fails.Add(failure.New(fails.ErrCritical, failure.Message("入稿前のイスが存在しています")))
		return
	}

	err = c.PostChairs(ctx, chairs)
	if err != nil {
		fails.Add(failure.Translate(err, fails.ErrCritical))
		return
	}

	chair, err = c.GetChairDetailFromID(ctx, id)
	if err != nil {
		fails.Add(failure.Translate(err, fails.ErrCritical))
		return
	}
	if chair == nil {
		fails.Add(failure.Translate(err, fails.ErrCritical, failure.Message("入稿したイスのデータが不正です")))
		return
	}
	if err := checkChairEqualToAsset(chair); err != nil {
		fails.Add(failure.Translate(err, fails.ErrCritical, failure.Message("入稿したイスのデータが不正です")))
		return
	}
}
