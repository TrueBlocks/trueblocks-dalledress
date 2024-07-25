package messages

import (
	"context"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/base"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type Message int

const (
	Completed Message = iota
	Error
	Warn
	Progress
	Server
)

type ProgressMsg struct {
	Address base.Address `json:"address"`
	Have    int64        `json:"have"`
	Want    int64        `json:"want"`
}

type ServerMsg struct {
	Name string `json:"name"`
}

func MessageType(msg Message) string {
	m := map[Message]string{
		Completed: "Completed",
		Error:     "Error",
		Warn:      "Warn",
		Progress:  "Progress",
		Server:    "Server",
	}
	return m[msg]
}

func SendMessage(ctx context.Context, addr base.Address, msg Message, data interface{}) {
	runtime.EventsEmit(ctx, MessageType(msg), data)
}
