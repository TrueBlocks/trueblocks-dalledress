package app

import (
	"strings"
	"sync"

	"github.com/TrueBlocks/trueblocks-core/sdk/v3"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/base"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/messages"
)

var ensLock sync.Mutex

func (a *App) ConvertToAddress(addr string) (base.Address, bool) {
	if !strings.HasSuffix(addr, ".eth") {
		ret := base.HexToAddress(addr)
		return ret, ret != base.ZeroAddr
	}

	ensLock.Lock()
	ensAddr, exists := a.ensMap[addr]
	ensLock.Unlock()

	if exists {
		return ensAddr, true
	}

	// Try to get an ENS or return the same input
	opts := sdk.NamesOptions{
		Terms: []string{addr},
	}
	if names, _, err := opts.Names(); err != nil {
		messages.Send(a.ctx, messages.Error, messages.NewErrorMsg(err))
		return base.ZeroAddr, false
	} else {
		if len(names) > 0 {
			ensLock.Lock()
			defer ensLock.Unlock()
			a.ensMap[addr] = names[0].Address
			return names[0].Address, true
		} else {
			ret := base.HexToAddress(addr)
			return ret, ret != base.ZeroAddr
		}
	}
}

func (a *App) AddrToName(addr base.Address) string {
	if name, exists := a.names.NamesMap[addr]; exists {
		return name.Name
	}
	return addr.Hex()
}
