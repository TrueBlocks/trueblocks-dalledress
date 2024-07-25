package app

import (
	"fmt"

	"github.com/TrueBlocks/trueblocks-dalledress/servers"
)

func (a *App) GetServer(name string) *servers.Server {
	switch name {
	case "scraper":
		return &a.Scraper.Server
	case "fileserver":
		return &a.FileServer.Server
	case "monitor":
		return &a.Monitor.Server
	case "ipfs":
		return &a.Ipfs.Server
	default:
		return nil
	}
}

func (a *App) ToggleServer(name string) error {
	switch name {
	case "scraper":
		return a.Scraper.Server.Toggle()
	case "fileserver":
		return a.FileServer.Toggle()
	case "monitor":
		return a.Monitor.Server.Toggle()
	case "ipfs":
		return a.Ipfs.Server.Toggle()
	default:
		return fmt.Errorf("server %s not found in ToggleServer", name)
	}
}

func (a *App) StartServers() {
	go a.Scraper.Run()
	// go a.FileServer.Run()
	go a.Monitor.Run()
	go a.Ipfs.Run()
}
