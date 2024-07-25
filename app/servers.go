package app

import (
	"fmt"

	"github.com/TrueBlocks/trueblocks-browse/servers"
)

func (a *App) GetServer(name string) *servers.Server {
	switch name {
	case "scraper":
		return &a.Scraper.Server
	default:
		return nil
	}
}

func (a *App) ToggleServer(name string) error {
	switch name {
	case "scraper":
		return a.Scraper.Server.Toggle()
	default:
		return fmt.Errorf("server %s not found in ToggleServer", name)
	}
}

func (a *App) StartServers() {
	go a.Scraper.Run()
}
