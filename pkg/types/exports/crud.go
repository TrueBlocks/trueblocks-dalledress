package exports

import (
	"fmt"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/crud"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/types"
)

type ExportsCrud struct{}

func NewExportsCrud() *ExportsCrud {
	return &ExportsCrud{}
}

func (c *ExportsCollection) Crud(
	dataFacet types.DataFacet,
	op crud.Operation,
	item interface{},
) error {
	// crudHandler := NewExportsCrud()
	return fmt.Errorf("operation %s not supported for exports", op)
}
