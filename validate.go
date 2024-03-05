package main

import (
	"strings"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/base"
)

func (a *App) validateInput(s string) (base.Address, bool) {
	if len(s) < 42 {
		return base.ZeroAddr, false
	}

	if strings.HasSuffix(s, ".eth") {
		if a.conn != nil {
			if addr, ok := a.conn.GetEnsAddress(s); ok {
				return base.HexToAddress(addr), true
			}
			return base.ZeroAddr, false
		}
	}

	return base.HexToAddress(s[:42]), base.IsValidAddress(s[:42])
}
