package score

import (
	"sync"

	"github.com/isucon10-qualify/isucon10-qualify/bench/parameter"
)

var (
	score     int64 = 0
	level     int64 = 0
	maxLevel  int64 = int64(len(parameter.BoundaryOfLevel)) - 1
	levelChan chan int64
	mu        sync.RWMutex
)

func init() {
	levelChan = make(chan int64, 1)
}

func IncrementScore() {
	mu.Lock()
	defer mu.Unlock()
	score++
	if level < maxLevel && score >= parameter.BoundaryOfLevel[level] {
		level++
		levelChan <- level
	}
}

func GetScore() int64 {
	mu.RLock()
	defer mu.RUnlock()
	return score
}

func GetLevel() int64 {
	mu.RLock()
	defer mu.RUnlock()
	return level
}

func LevelUp() chan int64 {
	return levelChan
}
