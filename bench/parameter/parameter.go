package parameter

import "time"

const (
	NumOfSearchChairInScenario     = 2
	NumOfSearchEstateInScenario    = 2
	NumOfCheckChairSearchPaging    = 2
	NumOfCheckEstateSearchPaging   = 2
	LimitOfChairSearchPageDepth    = 5
	LimitOfEstateSearchPageDepth   = 5
	NumOfCheckChairDetailPage      = 2
	NumOfCheckEstateDetailPage     = 2
	PerPageOfChairSearch           = 25
	PerPageOfEstateSearch          = 25
	MaxLengthOfNazotteResponse     = 50
	SleepTimeOnFailScenario        = 1500 * time.Millisecond
	SleepSwingOnFailScenario       = 500 // * time.Millisecond
	SleepTimeOnUserAway            = 500 * time.Millisecond
	SleepSwingOnUserAway           = 100 // * time.Millisecond
	SleepTimeOnBotInterval         = 500 * time.Millisecond
	SleepSwingOnBotInterval        = 100 // * time.Millisecond
	SleepBeforePostDraft           = 500 * time.Millisecond
	SleepSwingBeforePostDraft      = 100 // * time.Millisecond
	ThresholdTimeOfAbandonmentPage = 1000 * time.Millisecond
	DefaultAPITimeout              = 2000 * time.Millisecond
	InitializeTimeout              = 30 * time.Second
	VerifyTimeout                  = 10 * time.Second
	DraftTimeout                   = 5 * time.Second
	LoadTimeout                    = 60 * time.Second
)

var BoundaryOfLevel []int64 = []int64{
	300, 600, 800, 900, 1000,
	1100, 1200, 1300, 1450, 1600,
	1800, 2000,
}

type incWorkers struct {
	ChairSearchWorker         int
	EstateSearchWorker        int
	EstateNazotteSearchWorker int
	BotWorker                 int
	ChairDraftPostWorker      int
	EstateDraftPostWorker     int
}

// IncListOfWorkers 前のレベルとのWorkerの個数の差分を保持するList
var ListOfIncWorkers = []incWorkers{
	{ // level 00
		ChairSearchWorker:         3,
		EstateSearchWorker:        3,
		EstateNazotteSearchWorker: 0,
		BotWorker:                 0,
		ChairDraftPostWorker:      0,
		EstateDraftPostWorker:     0,
	},
	{ // level 01
		ChairSearchWorker:         0,
		EstateSearchWorker:        0,
		EstateNazotteSearchWorker: 3,
		BotWorker:                 0,
		ChairDraftPostWorker:      0,
		EstateDraftPostWorker:     0,
	},
	{ // level 02
		ChairSearchWorker:         0,
		EstateSearchWorker:        0,
		EstateNazotteSearchWorker: 0,
		BotWorker:                 5,
		ChairDraftPostWorker:      1,
		EstateDraftPostWorker:     1,
	},
	{ // level 03
		ChairSearchWorker:         1,
		EstateSearchWorker:        1,
		EstateNazotteSearchWorker: 1,
		BotWorker:                 1,
		ChairDraftPostWorker:      1,
		EstateDraftPostWorker:     1,
	},
	{ // level 04
		ChairSearchWorker:         1,
		EstateSearchWorker:        1,
		EstateNazotteSearchWorker: 1,
		BotWorker:                 1,
		ChairDraftPostWorker:      1,
		EstateDraftPostWorker:     1,
	},
	{ // level 05
		ChairSearchWorker:         1,
		EstateSearchWorker:        1,
		EstateNazotteSearchWorker: 1,
		BotWorker:                 1,
		ChairDraftPostWorker:      1,
		EstateDraftPostWorker:     1,
	},
	{ // level 06
		ChairSearchWorker:         1,
		EstateSearchWorker:        1,
		EstateNazotteSearchWorker: 1,
		BotWorker:                 1,
		ChairDraftPostWorker:      1,
		EstateDraftPostWorker:     1,
	},
	{ // level 07
		ChairSearchWorker:         1,
		EstateSearchWorker:        1,
		EstateNazotteSearchWorker: 1,
		BotWorker:                 1,
		ChairDraftPostWorker:      1,
		EstateDraftPostWorker:     1,
	},
	{ // level 08
		ChairSearchWorker:         1,
		EstateSearchWorker:        1,
		EstateNazotteSearchWorker: 1,
		BotWorker:                 1,
		ChairDraftPostWorker:      1,
		EstateDraftPostWorker:     1,
	},
	{ // level 09
		ChairSearchWorker:         1,
		EstateSearchWorker:        1,
		EstateNazotteSearchWorker: 1,
		BotWorker:                 1,
		ChairDraftPostWorker:      1,
		EstateDraftPostWorker:     1,
	},
	{ // level 10
		ChairSearchWorker:         1,
		EstateSearchWorker:        1,
		EstateNazotteSearchWorker: 1,
		BotWorker:                 1,
		ChairDraftPostWorker:      1,
		EstateDraftPostWorker:     1,
	},
	{ // level 11
		ChairSearchWorker:         1,
		EstateSearchWorker:        1,
		EstateNazotteSearchWorker: 1,
		BotWorker:                 1,
		ChairDraftPostWorker:      1,
		EstateDraftPostWorker:     1,
	},
}
