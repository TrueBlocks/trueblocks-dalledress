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
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/types/contracts"
	sdk "github.com/TrueBlocks/trueblocks-sdk/v5"
	// EXISTING_CODE
	// EXISTING_CODE
)

func (a *App) GetContractsPage(
	payload *types.Payload,
	first, pageSize int,
	sort sdk.SortSpec,
	filter string,
) (*contracts.ContractsPage, error) {
	collection := contracts.GetContractsCollection(payload)
	return getCollectionPage[*contracts.ContractsPage](collection, payload, first, pageSize, sort, filter)
}

func (a *App) GetContractsSummary(payload *types.Payload) types.Summary {
	collection := contracts.GetContractsCollection(payload)
	return collection.GetSummary()
}

func (a *App) ReloadContracts(payload *types.Payload) error {
	collection := contracts.GetContractsCollection(payload)
	collection.Reset(payload.DataFacet)
	collection.LoadData(payload.DataFacet)
	return nil
}

// EXISTING_CODE
func (a *App) GetContract(address string, abi *contracts.Abi) (*contracts.Contract, error) {
	collection := contracts.NewContractsCollection()
	return collection.GetContract(address, abi)
}

func (a *App) RefreshContract(contractState *contracts.Contract) error {
	collection := contracts.NewContractsCollection()
	return collection.RefreshContract(contractState)
}

// TODO: QUESTION - THIS SHOULD USE GETPAGE
// GetEvents retrieves events for a specific contract
func (a *App) GetEvents(payload *types.Payload, address string, count int) ([]*sdk.Log, error) {
	collection := contracts.GetContractsCollection(payload)
	return collection.GetEvents(address, count)
}

// EXISTING_CODE
