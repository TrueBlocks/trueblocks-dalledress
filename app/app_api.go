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
	"fmt"

	"github.com/TrueBlocks/trueblocks-dalledress/pkg/logging"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/types"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/types/abis"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/types/chunks"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/types/comparitoor"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/types/contracts"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/types/dalledress"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/types/exports"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/types/monitors"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/types/names"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/types/status"
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
	case "dalledress":
		return a.ReloadDalleDress(payload)
	case "comparitoor":
		return a.ReloadComparitoor(payload)
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
		"dalledress",
		"comparitoor",
	}
}

func getCollection(payload *types.Payload) types.Collection {
	switch payload.Collection {
	case "exports":
		return exports.GetExportsCollection(payload)
	case "monitors":
		return monitors.GetMonitorsCollection(payload)
	case "abis":
		return abis.GetAbisCollection(payload)
	case "names":
		return names.GetNamesCollection(payload)
	case "chunks":
		return chunks.GetChunksCollection(payload)
	case "contracts":
		return contracts.GetContractsCollection(payload)
	case "status":
		return status.GetStatusCollection(payload)
	case "dalledress":
		return dalledress.GetDalleDressCollection(payload)
	case "comparitoor":
		return comparitoor.GetComparitoorCollection(payload)
	default:
		logging.LogBackend(fmt.Sprintf("Warning: Unknown collection type: %s", payload.Collection))
		return nil
	}
}

// EXISTING_CODE
// EXISTING_CODE
