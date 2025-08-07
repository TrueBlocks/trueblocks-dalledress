// Copyright 2016, 2026 The Authors. All rights reserved.
// Use of this source code is governed by a license that can
// be found in the LICENSE file.
/*
 * Parts of this file were auto generated. Edit only those parts of
 * the code inside of 'EXISTING_CODE' tags.
 */

package app

// EXISTING_CODE
import (
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/types"
)

// EXISTING_CODE

// Reload dispatches reload requests to the appropriate view-specific reload function
func (a *App) Reload(payload *types.Payload) error {
	switch a.GetLastView() {
	case "exports":
		return a.ReloadExports(payload)
	case "monitors":
		return a.ReloadMonitors(payload)
	case "abis":
		return a.ReloadAbis(payload)
	case "names":
		return a.ReloadNames(payload)
	case "chunks":
		return a.ReloadChunks(payload)
	case "contracts":
		return a.ReloadContracts(payload)
	case "status":
		return a.ReloadStatus(payload)
	default:
		panic("unknown view in Reload" + a.GetLastView())
	}
}

// GetRegisteredViews returns all registered view names
func (a *App) GetRegisteredViews() []string {
	return []string{
		"exports",
		"monitors",
		"abis",
		"names",
		"chunks",
		"contracts",
		"status",
	}
}

// EXISTING_CODE
// EXISTING_CODE
