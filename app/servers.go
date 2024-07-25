package app

import (
	"time"

	"github.com/TrueBlocks/trueblocks-browse/servers"
)

var ss = servers.Server{
	Name:    "Testing Server",
	State:   servers.Running,
	Sleep:   1,
	Started: time.Now(),
	Runs:    12,
}

func (a *App) GetServer(name string) *servers.Server {
	return &ss
}
