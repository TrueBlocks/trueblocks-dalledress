package servers

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/base"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/colors"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/messages"
)

type Server struct {
	Name    string          `json:"name"`
	Sleep   time.Duration   `json:"sleep"`
	Started time.Time       `json:"started"`
	Runs    int             `json:"runs"`
	State   State           `json:"state"`
	MsgCtx  context.Context `json:"-"`
}

func (s *Server) Run() error {
	s.State = Running
	s.Notify("Run")
	return nil
}

func (s *Server) Stop() error {
	s.State = Stopped
	s.Notify("Stopped")
	return nil
}

func (s *Server) Pause() error {
	s.State = Paused
	s.Notify("Paused")
	return nil
}

func (s *Server) Toggle() error {
	if s.State == Running {
		return s.Pause()
	}
	return s.Run()
}

func (s *Server) Tick() int {
	s.Runs++
	return s.Runs
}

func (s *Server) getColor() string {
	return colors.Yellow
}

func (s *Server) Notify(msg ...string) {
	s.Tick()
	fmt.Printf("Notify %s%-10.10s (% 5d-% 5.2f): %s%s\n",
		s.getColor(),
		s.Name,
		s.Runs,
		float64(time.Since(s.Started))/float64(time.Second),
		msg,
		colors.Off,
	)
	messages.SendMessage(s.MsgCtx, base.ZeroAddr, messages.Server, messages.ServerMsg{
		Name: strings.ToLower(s.Name),
	})
}

type Serverer *interface {
	Run() error
	Stop() error
	Pause() error
	Tick() int
}
