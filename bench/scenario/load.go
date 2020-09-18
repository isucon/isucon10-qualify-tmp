package scenario

import (
	"context"
	"log"
	"math/rand"
	"time"

	"github.com/isucon10-qualify/isucon10-qualify/bench/asset"
	"github.com/isucon10-qualify/isucon10-qualify/bench/client"
	"github.com/isucon10-qualify/isucon10-qualify/bench/fails"
	"github.com/isucon10-qualify/isucon10-qualify/bench/parameter"
	"github.com/isucon10-qualify/isucon10-qualify/bench/score"
	"github.com/morikuni/failure"
)

func runEstateSearchWorker(ctx context.Context) {

	c := client.NewClient(false)

	for {
		r := rand.Intn(100)
		t := time.NewTimer(time.Duration(r) * time.Millisecond)
		select {
		case <-t.C:
		case <-ctx.Done():
			t.Stop()
			return
		}
		err := estateSearchScenario(ctx, c)
		if err != nil {
			code, _ := failure.CodeOf(err)
			if code == fails.ErrTimeout {
				r := rand.Intn(parameter.SleepSwingOnUserAway) - parameter.SleepSwingOnUserAway*0.5
				s := parameter.SleepTimeOnFailScenario + time.Duration(r)*time.Millisecond
				t = time.NewTimer(s)
			} else {
				r := rand.Intn(parameter.SleepSwingOnFailScenario) - parameter.SleepSwingOnFailScenario*0.5
				s := parameter.SleepTimeOnFailScenario + time.Duration(r)*time.Millisecond
				t = time.NewTimer(s)
			}
			select {
			case <-t.C:
			case <-ctx.Done():
				t.Stop()
				return
			}
		}
	}
}

func runChairSearchWorker(ctx context.Context) {

	c := client.NewClient(false)

	for {
		r := rand.Intn(100)
		t := time.NewTimer(time.Duration(r) * time.Millisecond)
		select {
		case <-t.C:
		case <-ctx.Done():
			t.Stop()
			return
		}
		err := chairSearchScenario(ctx, c)
		if err != nil {
			code, _ := failure.CodeOf(err)
			if code == fails.ErrTimeout {
				r := rand.Intn(parameter.SleepSwingOnUserAway) - parameter.SleepSwingOnUserAway*0.5
				s := parameter.SleepTimeOnFailScenario + time.Duration(r)*time.Millisecond
				t = time.NewTimer(s)
			} else {
				r := rand.Intn(parameter.SleepSwingOnFailScenario) - parameter.SleepSwingOnFailScenario*0.5
				s := parameter.SleepTimeOnFailScenario + time.Duration(r)*time.Millisecond
				t = time.NewTimer(s)
			}
			select {
			case <-t.C:
			case <-ctx.Done():
				t.Stop()
				return
			}
		}
	}
}

func runEstateNazotteSearchWorker(ctx context.Context) {

	c := client.NewClient(false)

	for {
		r := rand.Intn(100)
		t := time.NewTimer(time.Duration(r) * time.Millisecond)
		select {
		case <-t.C:
		case <-ctx.Done():
			t.Stop()
			return
		}
		err := estateNazotteSearchScenario(ctx, c)
		if err != nil {
			code, _ := failure.CodeOf(err)
			if code == fails.ErrTimeout {
				r := rand.Intn(parameter.SleepSwingOnUserAway) - parameter.SleepSwingOnUserAway*0.5
				s := parameter.SleepTimeOnFailScenario + time.Duration(r)*time.Millisecond
				t = time.NewTimer(s)
			} else {
				r := rand.Intn(parameter.SleepSwingOnFailScenario) - parameter.SleepSwingOnFailScenario*0.5
				s := parameter.SleepTimeOnFailScenario + time.Duration(r)*time.Millisecond
				t = time.NewTimer(s)
			}
			select {
			case <-t.C:
			case <-ctx.Done():
				t.Stop()
				return
			}
		}
	}
}

func runBotWorker(ctx context.Context) {

	c := client.NewClient(true)

	for {
		go botScenario(ctx, c)
		r := rand.Intn(parameter.SleepSwingOnBotInterval) - parameter.SleepSwingOnBotInterval*0.5
		s := parameter.SleepTimeOnBotInterval + time.Duration(r)*time.Millisecond
		t := time.NewTimer(s)
		select {
		case <-t.C:
		case <-ctx.Done():
			t.Stop()
			return
		}
	}
}

func runChairDraftPostWorker(ctx context.Context) {
	c := client.NewClientForDraft()

	r := rand.Intn(parameter.SleepSwingBeforePostDraft) - parameter.SleepSwingBeforePostDraft*0.5
	s := parameter.SleepBeforePostDraft + time.Duration(r)*time.Millisecond
	t := time.NewTimer(s)
	select {
	case <-t.C:
		filePath, err := asset.ChairDraftFiles.Next()
		if err != nil {
			return
		}
		chairDraftPostScenario(ctx, c, filePath)
	case <-ctx.Done():
		t.Stop()
		return
	}
}

func runEstateDraftPostWorker(ctx context.Context) {
	c := client.NewClientForDraft()

	r := rand.Intn(parameter.SleepSwingBeforePostDraft) - parameter.SleepSwingBeforePostDraft*0.5
	s := parameter.SleepBeforePostDraft + time.Duration(r)*time.Millisecond
	t := time.NewTimer(s)
	select {
	case <-t.C:
		filePath, err := asset.EstateDraftFiles.Next()
		if err != nil {
			return
		}
		estateDraftPostScenario(ctx, c, filePath)
	case <-ctx.Done():
		t.Stop()
	}
}

func checkWorkers(ctx context.Context) {
	for {
		select {
		case level := <-score.LevelUp():
			log.Println("負荷レベルが上昇しました。")
			incWorkers := parameter.ListOfIncWorkers[level]
			for i := 0; i < incWorkers.ChairSearchWorker; i++ {
				go runChairSearchWorker(ctx)
			}
			for i := 0; i < incWorkers.EstateSearchWorker; i++ {
				go runEstateSearchWorker(ctx)
			}
			for i := 0; i < incWorkers.EstateNazotteSearchWorker; i++ {
				go runEstateNazotteSearchWorker(ctx)
			}
			for i := 0; i < incWorkers.BotWorker; i++ {
				go runBotWorker(ctx)
			}
			for i := 0; i < incWorkers.ChairDraftPostWorker; i++ {
				go runChairDraftPostWorker(ctx)
			}
			for i := 0; i < incWorkers.EstateDraftPostWorker; i++ {
				go runEstateDraftPostWorker(ctx)
			}
		case <-ctx.Done():
			return
		}
	}
}

func Load(ctx context.Context) {
	level := score.GetLevel()
	incWorkers := parameter.ListOfIncWorkers[level]

	// 物件検索をして、資料請求をするシナリオ
	for i := 0; i < incWorkers.ChairSearchWorker; i++ {
		go runChairSearchWorker(ctx)
	}

	// イス検索から物件ページに行き、資料請求をするまでのシナリオ
	for i := 0; i < incWorkers.EstateSearchWorker; i++ {
		go runEstateSearchWorker(ctx)
	}

	// なぞって検索をするシナリオ
	for i := 0; i < incWorkers.EstateNazotteSearchWorker; i++ {
		go runEstateNazotteSearchWorker(ctx)
	}

	// ボットによる検索シナリオ
	for i := 0; i < incWorkers.BotWorker; i++ {
		go runBotWorker(ctx)
	}

	// イスの入稿シナリオ
	for i := 0; i < incWorkers.ChairDraftPostWorker; i++ {
		go runChairDraftPostWorker(ctx)
	}

	// 物件の入稿シナリオ
	for i := 0; i < incWorkers.EstateDraftPostWorker; i++ {
		go runEstateDraftPostWorker(ctx)
	}

	go checkWorkers(ctx)
}
