// Copyright 2016, 2026 The TrueBlocks Authors. All rights reserved.
// Use of this source code is governed by a license that can
// be found in the LICENSE file.
/*
 * Parts of this file were auto generated. Edit only those parts of
 * the code inside of 'EXISTING_CODE' tags.
 */

package app

import (
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/types"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/types/chunks"
	sdk "github.com/TrueBlocks/trueblocks-sdk/v5"
	// EXISTING_CODE
	// EXISTING_CODE
)

func (a *App) GetChunksPage(
	payload *types.Payload,
	first, pageSize int,
	sort sdk.SortSpec,
	filter string,
) (*chunks.ChunksPage, error) {
	collection := chunks.GetChunksCollection(payload)
	return getCollectionPage[*chunks.ChunksPage](collection, payload, first, pageSize, sort, filter)
}

func (a *App) GetChunksSummary(payload *types.Payload) types.Summary {
	collection := chunks.GetChunksCollection(payload)
	return collection.GetSummary()
}

func (a *App) ReloadChunks(payload *types.Payload) error {
	collection := chunks.GetChunksCollection(payload)
	collection.Reset(payload.DataFacet)
	collection.LoadData(payload.DataFacet)
	return nil
}

// EXISTING_CODE
// EXISTING_CODE
