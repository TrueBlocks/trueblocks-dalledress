// Copyright 2016, 2026 The Authors. All rights reserved.
// Use of this source code is governed by a license that can
// be found in the LICENSE file.
/*
 * Parts of this file were auto generated. Edit only those parts of
 * the code inside of 'EXISTING_CODE' tags.
 */

package app

import (
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/types"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/types/dalledress"

	//
	sdk "github.com/TrueBlocks/trueblocks-sdk/v5"
	// EXISTING_CODE
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/crud"
	// EXISTING_CODE
)

func (a *App) GetDalleDressPage(
	payload *types.Payload,
	first, pageSize int,
	sort sdk.SortSpec,
	filter string,
) (*dalledress.DalleDressPage, error) {
	collection := dalledress.GetDalleDressCollection(payload)
	return getCollectionPage[*dalledress.DalleDressPage](collection, payload, first, pageSize, sort, filter)
}

func (a *App) GetDalleDressSummary(payload *types.Payload) types.Summary {
	collection := dalledress.GetDalleDressCollection(payload)
	return collection.GetSummary()
}

func (a *App) ReloadDalleDress(payload *types.Payload) error {
	collection := dalledress.GetDalleDressCollection(payload)
	collection.Reset(payload.DataFacet)
	collection.LoadData(payload.DataFacet)
	return nil
}

// GetDalleDressConfig returns the view configuration for dalledress
func (a *App) GetDalleDressConfig(payload types.Payload) (*types.ViewConfig, error) {
	collection := dalledress.GetDalleDressCollection(&payload)
	return collection.GetConfig()
}

// EXISTING_CODE
func (a *App) SeriesCrud(
	payload *types.Payload,
	op crud.Operation,
	item *dalledress.Series,
) error {
	collection := dalledress.GetDalleDressCollection(payload)
	return collection.SeriesCrud(payload, op, item)
}

// EXISTING_CODE
