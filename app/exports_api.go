// ADD_ROUTE
package app

import (
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/crud"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/types"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/types/exports"
	sdk "github.com/TrueBlocks/trueblocks-sdk/v5"
)

func (a *App) GetExportsPage(
	dataFacet types.DataFacet,
	first, pageSize int,
	sort sdk.SortSpec,
	filter string,
	chain string,
	address string,
) (*exports.ExportsPage, error) {
	collection := exports.GetExportsCollection(chain, address)
	return getCollectionPage[*exports.ExportsPage](collection, dataFacet, first, pageSize, sort, filter)
}

func (a *App) ExportsCrud(
	dataFacet types.DataFacet,
	op crud.Operation,
	item interface{},
	chain string,
	address string,
) error {
	collection := exports.GetExportsCollection(chain, address)
	return collection.Crud(dataFacet, op, item)
}

func (a *App) GetExportsCount(
	dataFacet types.DataFacet,
	chain string,
	address string,
) (uint64, error) {
	return exports.GetExportsCount(chain, address, string(dataFacet))
}

func (a *App) LoadExportsData(
	dataFacet types.DataFacet,
	chain string,
	address string,
) {
	collection := exports.GetExportsCollection(chain, address)
	collection.LoadData(dataFacet)
}

func (a *App) ResetExportsData(
	dataFacet types.DataFacet,
	chain string,
	address string,
) {
	collection := exports.GetExportsCollection(chain, address)
	collection.Reset(dataFacet)
}

func (a *App) ExportsNeedsUpdate(
	dataFacet types.DataFacet,
	chain string,
	address string,
) bool {
	collection := exports.GetExportsCollection(chain, address)
	return collection.NeedsUpdate(dataFacet)
}

func (a *App) GetExportsSummary(chain string, address string) types.Summary {
	collection := exports.GetExportsCollection(chain, address)
	return collection.GetSummary()
}

// ADD_ROUTE
