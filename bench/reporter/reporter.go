package reporter

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"sync"

	"github.com/isucon10-qualify/isucon10-qualify/bench/score"
)

type Stdout struct {
	Pass     bool      `json:"pass"`
	Score    int64     `json:"score"`
	Messages []Message `json:"messages"`
	Reason   string    `json:"reason"`
	Language string    `json:"language"`
}

var mu sync.RWMutex
var writer io.Writer
var stdout *Stdout

func init() {
	stdout = &Stdout{
		Pass:     false,
		Score:    0,
		Messages: make([]Message, 0),
		Reason:   "",
		Language: "",
	}
}

func Report(msgs []string, critical, application, trivial int) error {
	err := update(msgs, critical, application, trivial)
	if err != nil {
		return err
	}

	mu.RLock()
	defer mu.RUnlock()
	bytes, err := json.Marshal(stdout)
	if err != nil {
		return err
	}
	fmt.Println(string(bytes))
	return nil
}

func SetPassed(p bool) {
	mu.Lock()
	defer mu.Unlock()
	stdout.Pass = p
}

func SetReason(reason string) {
	mu.Lock()
	defer mu.Unlock()
	stdout.Reason = reason
}

func update(msgs []string, critical, application, trivial int) error {
	mu.Lock()
	defer mu.Unlock()

	row := score.GetScore()
	deducation := int64(application * 50)
	score := row - deducation
	if score < 0 {
		stdout.Pass = false
		stdout.Reason = "スコアが0点を下回りました"
		score = 0
	}
	stdout.Score = score

	stdout.Messages = UniqMsgs(msgs)
	return nil
}

func SetLanguage(language string) {
	mu.Lock()
	defer mu.Unlock()
	stdout.Language = language
}

type Message struct {
	Text  string `json:"text"`
	Count int    `json:"count"`
}

func UniqMsgs(allMsgs []string) []Message {
	if len(allMsgs) == 0 {
		return []Message{}
	}

	sort.Strings(allMsgs)
	msgs := make([]Message, 0, len(allMsgs))

	preMsg := allMsgs[0]
	cnt := 0

	// 適当にuniqする
	for _, msg := range allMsgs {
		if preMsg != msg {
			msgs = append(msgs, Message{
				Text:  preMsg,
				Count: cnt,
			})
			preMsg = msg
			cnt = 1
		} else {
			cnt++
		}
	}
	msgs = append(msgs, Message{
		Text:  preMsg,
		Count: cnt,
	})

	return msgs
}
