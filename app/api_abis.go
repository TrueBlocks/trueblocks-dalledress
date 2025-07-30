package app

import (
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/crud"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/types"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/types/abis"
	sdk "github.com/TrueBlocks/trueblocks-sdk/v5"
)

func (a *App) GetAbisPage(
	payload *types.Payload,
	first, pageSize int,
	sort sdk.SortSpec,
	filter string,
) (*abis.AbisPage, error) {
	collection := abis.GetAbisCollection(payload)
	return getCollectionPage[*abis.AbisPage](collection, payload, first, pageSize, sort, filter)
}

func (a *App) AbisCrud(
	payload *types.Payload,
	op crud.Operation,
	item *abis.Abi,
) error {
	collection := abis.GetAbisCollection(payload)
	return collection.Crud(payload, op, item)
}

func (a *App) GetAbisSummary(payload *types.Payload) types.Summary {
	collection := abis.GetAbisCollection(payload)
	return collection.GetSummary()
}

func (a *App) ReloadAbis(payload *types.Payload) error {
	collection := abis.GetAbisCollection(payload)
	collection.Reset(payload.DataFacet)
	collection.LoadData(payload.DataFacet)
	return nil
}
