// ADD_ROUTE
package app

import (
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/crud"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/types"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/types/exports"
	sdk "github.com/TrueBlocks/trueblocks-sdk/v5"
)

func (a *App) GetExportsPage(
	kind types.ListKind,
	first, pageSize int,
	sort sdk.SortSpec,
	filter string,
	chain string,
	address string,
) (*exports.ExportsPage, error) {
	collection := exports.GetExportsCollection(chain, address)
	return getCollectionPage[*exports.ExportsPage](collection, kind, first, pageSize, sort, filter)
}

func (a *App) ExportsCrud(
	kind types.ListKind,
	op crud.Operation,
	item interface{},
	chain string,
	address string,
) error {
	collection := exports.GetExportsCollection(chain, address)
	return collection.Crud(kind, op, item)
}

func (a *App) GetExportsCount(
	kind types.ListKind,
	chain string,
	address string,
) (uint64, error) {
	return exports.GetExportsCount(chain, address, string(kind))
}

func (a *App) LoadExportsData(
	kind types.ListKind,
	chain string,
	address string,
) {
	collection := exports.GetExportsCollection(chain, address)
	collection.LoadData(kind)
}

func (a *App) ResetExportsData(
	kind types.ListKind,
	chain string,
	address string,
) {
	collection := exports.GetExportsCollection(chain, address)
	collection.Reset(kind)
}

func (a *App) ExportsNeedsUpdate(
	kind types.ListKind,
	chain string,
	address string,
) bool {
	collection := exports.GetExportsCollection(chain, address)
	return collection.NeedsUpdate(kind)
}

func (a *App) GetExportsSummary(chain string, address string) types.Summary {
	collection := exports.GetExportsCollection(chain, address)
	return collection.GetSummary()
}

// ADD_ROUTE
