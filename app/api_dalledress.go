// Copyright 2016, 2026 The Authors. All rights reserved.
// Use of this source code is governed by a license that can
// be found in the LICENSE file.
/*
 * Parts of this file were auto generated. Edit only those parts of
 * the code inside of 'EXISTING_CODE' tags.
 */

package app

import (
	"fmt"

	"github.com/TrueBlocks/trueblocks-dalledress/pkg/types"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/types/dalledress"

	//
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/crud"
	sdk "github.com/TrueBlocks/trueblocks-sdk/v5"

	// EXISTING_CODE
	dalle "github.com/TrueBlocks/trueblocks-dalle/v2"
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

func (a *App) DalleDressCrud(
	payload *types.Payload,
	op crud.Operation,
	item *any,
) error {
	collection := dalledress.GetDalleDressCollection(payload)
	return collection.Crud(payload, op, item)
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
func (a *App) Speak(payload *types.Payload, series string) (string, error) {
	if payload == nil || payload.Address == "" {
		return "", nil
	}
	if series == "" {
		series = "empty"
	}
	return dalle.Speak(series, payload.Address)
}

func (a *App) ReadToMe(payload *types.Payload, series string) (string, error) {
	if payload == nil || payload.Address == "" {
		return "", nil
	}
	if series == "" {
		series = "empty"
	}
	return dalle.ReadToMe(series, payload.Address)
}

func (a *App) GetDalleAudioURL(payload *types.Payload, series string) (string, error) {
	if payload == nil || payload.Address == "" {
		return "", nil
	}
	if series == "" {
		series = "empty"
	}
	base := ""
	if a.fileServer != nil {
		base = a.fileServer.GetBaseURL()
	}
	return dalle.AudioURL(base, series, payload.Address)
}

func (a *App) FromTemplate(payload *types.Payload, templateStr string) (string, error) {
	if payload == nil || payload.Address == "" {
		return "", fmt.Errorf("address is required")
	}
	if templateStr == "" {
		return "", fmt.Errorf("template string is required")
	}

	dd, err := a.Dalle.MakeDalleDress(payload.Address)
	if err != nil {
		return "", fmt.Errorf("failed to create DalleDress: %w", err)
	}

	return dd.FromTemplate(templateStr)
}

// EXISTING_CODE
