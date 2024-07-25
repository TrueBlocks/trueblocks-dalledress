package servers

import (
	"time"

	// "github.com/TrueBlocks/trueblocks-core/sdk/v3"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/logger"
)

type Scraper struct {
	Server `json:"server"`
}

func NewScraper(name string, sleep time.Duration) *Scraper {
	return &Scraper{
		Server: Server{
			Name:    name,
			Sleep:   sleep,
			Color:   "yellow",
			State:   Paused,
			Started: time.Now(),
		},
	}
}

func (s *Scraper) Run() {
	logger.Info("Starting scraper...")

	for {
		if s.Server.State == Running {
			// opts := sdk.ScrapeOptions{
			// 	RunCount: 1,
			// }
			// _ = opts.Scrape()
			s.Server.Notify()
		}
		time.Sleep(s.Sleep * time.Millisecond)
	}
}

func (s *Scraper) Stop() error {
	return s.Server.Stop()
}

func (s *Scraper) Pause() error {
	return s.Server.Pause()
}

func (s *Scraper) Toggle() error {
	return s.Server.Toggle()
}

func (s *Scraper) Tick() int {
	return s.Server.Tick()
}
