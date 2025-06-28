// Copyright 2016, 2026 The TrueBlocks Authors. All rights reserved.
// Use of this source code is governed by a license that can
// be found in the LICENSE file.
/*
 * Parts of this file were auto generated. Edit only those parts of
 * the code inside of 'EXISTING_CODE' tags.
 */

package app

import (
	// EXISTING_CODE

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/base"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/crud"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/types"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/types/names"
	sdk "github.com/TrueBlocks/trueblocks-sdk/v5"
	// EXISTING_CODE
)

func (a *App) GetNamesPage(
	payload *types.Payload,
	first, pageSize int,
	sort sdk.SortSpec,
	filter string,
) (*names.NamesPage, error) {
	collection := names.GetNamesCollection(payload)
	return getCollectionPage[*names.NamesPage](collection, payload, first, pageSize, sort, filter)
}

func (a *App) NamesCrud(
	payload *types.Payload,
	op crud.Operation,
	item *names.Name,
) error {
	collection := names.GetNamesCollection(payload)
	return collection.Crud(payload, op, item)
}

func (a *App) GetNamesSummary(payload *types.Payload) types.Summary {
	collection := names.GetNamesCollection(payload)
	return collection.GetSummary()
}

func (a *App) ReloadNames(payload *types.Payload) error {
	collection := names.GetNamesCollection(payload)
	collection.Reset(payload.DataFacet)
	collection.LoadData(payload.DataFacet)
	return nil
}

// EXISTING_CODE
func (a *App) NameFromAddress(address string) (*names.Name, bool) {
	collection := names.GetNamesCollection(&types.Payload{})
	return collection.NameFromAddress(base.HexToAddress(address))
}

// EXISTING_CODE
