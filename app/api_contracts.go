package app

import (
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/types"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/types/contracts"
	sdk "github.com/TrueBlocks/trueblocks-sdk/v5"
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
