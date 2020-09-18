package main

import (
	"context"
	"flag"
	"log"
	"math/rand"
	"net"
	"os"
	"time"

	"github.com/isucon10-qualify/isucon10-qualify/bench/asset"
	"github.com/isucon10-qualify/isucon10-qualify/bench/client"
	"github.com/isucon10-qualify/isucon10-qualify/bench/fails"
	"github.com/isucon10-qualify/isucon10-qualify/bench/reporter"
	"github.com/isucon10-qualify/isucon10-qualify/bench/scenario"
	"github.com/isucon10-qualify/isucon10-qualify/bench/score"
	"github.com/morikuni/failure"
)

type Config struct {
	TargetURLStr string
	TargetHost   string

	AllowedIPs []net.IP
}

func init() {
	rand.Seed(time.Now().UnixNano())
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

func main() {
	defer func() {
		reporter.Report(fails.Get())
	}()

	defer func() {
		err := recover()
		if err, ok := err.(error); ok {
			err = failure.Translate(err, fails.ErrBenchmarker)
			fails.Add(err)
		}
	}()

	flags := flag.NewFlagSet("isucon10-qualify", flag.ContinueOnError)
	flags.SetOutput(os.Stderr)

	conf := Config{}
	dataDir := ""
	fixtureDir := ""

	flags.StringVar(&conf.TargetURLStr, "target-url", "http://localhost:1323", "target url")
	flags.StringVar(&dataDir, "data-dir", "../initial-data", "data directory")
	flags.StringVar(&fixtureDir, "fixture-dir", "../webapp/fixture", "fixture directory")

	err := flags.Parse(os.Args[1:])
	if err != nil {
		err = failure.Translate(err, fails.ErrBenchmarker, failure.Message("コマンドライン引数のパースに失敗しました"))
		fails.Add(err)
		reporter.SetPassed(false)
		reporter.SetReason("コマンドライン引数のパースに失敗しました")
		return
	}

	err = client.SetShareTargetURLs(
		conf.TargetURLStr,
		conf.TargetHost,
	)
	if err != nil {
		fails.Add(failure.Translate(err, fails.ErrBenchmarker))
		reporter.SetPassed(false)
		reporter.SetReason("ベンチ対象サーバーのURLが不正です")
		return
	}

	asset.Initialize(context.Background(), dataDir, fixtureDir)
	msgs := fails.GetMsgs()
	if len(msgs) > 0 {
		log.Println("asset initialize failed")
		reporter.SetPassed(false)
		reporter.SetReason("ベンチマーカーの初期化に失敗しました")
		return
	}

	log.Println("=== initialize ===")
	initRes := scenario.Initialize(context.Background())
	msgs = fails.GetMsgs()
	if len(msgs) > 0 {
		log.Println("initialize failed")
		reporter.SetPassed(false)
		reporter.SetReason("POST /initializeに失敗しました")
		return
	}

	reporter.SetLanguage(initRes.Language)

	log.Println("=== verify ===")
	scenario.Verify(context.Background(), dataDir, fixtureDir)
	msgs = fails.GetMsgs()
	if len(msgs) > 0 {
		log.Println("verify failed")
		reporter.SetPassed(false)
		reporter.SetReason("アプリケーション互換性チェックに失敗しました")
		return
	}

	log.Println("=== validation ===")
	scenario.Validation(context.Background())
	log.Printf("最終的な負荷レベル: %d", score.GetLevel())

	// ベンチマーク終了時にcritical errorが1つ以上、もしくはapplication errorが10回以上で失格
	msgs, critical, application, _ := fails.Get()
	isPassed := true

	if critical > 0 {
		isPassed = false
		reporter.SetReason("致命的なエラーが発生しました")
	} else if application >= 10 {
		isPassed = false
		reporter.SetReason("アプリケーションエラーが10回以上発生しました")
	} else {
		reporter.SetReason("OK")
	}

	reporter.SetPassed(isPassed)
}
