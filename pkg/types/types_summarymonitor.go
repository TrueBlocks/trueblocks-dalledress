package types

import (
	"encoding/json"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/base"
	coreTypes "github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/types"
)

type SummaryMonitor struct {
	coreTypes.Monitor
	NMonitors  int64                              `json:"nMonitors"`
	NNamed     int64                              `json:"nNamed"`
	NDeleted   int64                              `json:"nDeleted"`
	MonitorMap map[base.Address]coreTypes.Monitor `json:"monitorMap"`
	Monitors   []coreTypes.Monitor                `json:"monitors"`
}

func (s SummaryMonitor) String() string {
	bytes, _ := json.Marshal(s)
	return string(bytes)
}

func (s *SummaryMonitor) ShallowCopy() SummaryMonitor {
	return SummaryMonitor{
		Monitor:   s.Monitor,
		NNamed:    s.NNamed,
		NDeleted:  s.NDeleted,
		NMonitors: s.NMonitors,
	}
}

func (s *SummaryMonitor) Summarize() {
	for _, mon := range s.Monitors {
		s.NMonitors++
		if mon.Deleted {
			s.NDeleted++
		}
		if len(mon.Name) > 0 {
			s.NNamed++
		}
		s.FileSize += mon.FileSize
		s.NRecords += mon.NRecords
	}
}
