package coordinator

import (
	"context"
	"github.com/pkg/errors"
	"net/http"
	"net/url"
)

type VideoState int

var (
	NoProxiesAvailableErr = errors.New("No proxies available")
)

const (
	TX_IDLE       VideoState = 0
	TX_RESET                 = 1
	TX_AWAY                  = 2
	TX_CONFIRMING            = 11
	TX_PVT                   = 12
	TX_GRP                   = 13
	TX_KILL_MODEL            = 15
	RX_IDLE                  = 90
	RX_PVT                   = 91
	RX_VOY                   = 92
	RX_GRP                   = 93
	OFFLINE                  = 127
)

var IgnoredStates = map[VideoState]bool{
	TX_PVT:  true,
	TX_GRP:  true,
	TX_AWAY: true,
	RX_IDLE: true,
	OFFLINE: true,
}

type AccessLevel int

const (
	GUEST   AccessLevel = 0
	BASIC               = 1
	PREMIUM             = 2
	MODEL               = 4
	ADMIN               = 5
)

type MyFreeCamsPerformer struct {
	AccessLevel AccessLevel `json:"lv"`
	StageName   string      `json:"nm"`
	Pid         int64       `json:"pid"`
	SessionId   int64       `json:"sid"`
	UserId      int64       `json:"uid"`
	VideoState  VideoState  `json:"vs"`
	ViewerCount int64       `json:"viewerCount"`

	User struct {
		Age       int    `json:"age"`
		Avatar    int    `json:"avatar"`
		Blurb     string `json:"blurb"`
		Camserv   int    `json:"camserv"`
		ChatBg    int    `json:"chat_bg"`
		ChatColor string `json:"chat_color"`
		ChatFont  int    `json:"chat_font"`
		ChatOpt   int    `json:"chat_opt"`
		City      string `json:"city"`
		Country   string `json:"country"`
		Creation  int    `json:"creation"`
		Ethnic    string `json:"ethnic"`
		Photos    int    `json:"photos"`
		Profile   int    `json:"profile"`
	} `json:"u"`

	Meta struct {
		Camscore  float64 `json:"camscore"`
		Continent string  `json:"continent"`
		Flags     int     `json:"flags"`
		Kbit      int     `json:"kbit"`
		Lastnews  int     `json:"lastnews"`
		Mg        int     `json:"mg"`
		Missmfc   int     `json:"missmfc"`
		NewModel  int     `json:"new_model"`
		Rank      int     `json:"rank"`
		Rc        int     `json:"rc"`
		Topic     string  `json:"topic"`
	} `json:"m"`
}

func (performer *MyFreeCamsPerformer) IsComplete() bool {
	if performer.StageName == "" {
		return false
	}

	if performer.ViewerCount == 0 {
		return false
	}

	if _, ok := IgnoredStates[performer.VideoState]; ok {
		return false
	}

	if performer.AccessLevel != MODEL {
		return false
	}

	if performer.User.Camserv == 0 {
		return false
	}

	return true
}

type Scraper interface {
	Scrape(ctx context.Context) error
}

type MyfreeCamsScraper interface {
	Scrape(ctx context.Context) error
	GetPerformers() ([]*MyFreeCamsPerformer, int)
	GetSessionId() int64
}

type ScraperProxyService interface {
	Run(ctx context.Context)
	GetProxy(r *http.Request) (*url.URL, error)
}
