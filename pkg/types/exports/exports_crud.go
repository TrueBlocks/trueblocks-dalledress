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

func (c *ExportsCrud) Create(item interface{}) error {
	return types.NewValidationError("exports", "", "create",
		fmt.Errorf("exports collection is read-only"))
}

func (c *ExportsCrud) Update(item interface{}) error {
	return types.NewValidationError("exports", "", "update",
		fmt.Errorf("exports collection is read-only"))
}

func (c *ExportsCrud) Delete(item interface{}) error {
	return types.NewValidationError("exports", "", "delete",
		fmt.Errorf("exports collection is read-only"))
}

func (ec *ExportsCollection) Crud(
	dataFacet types.DataFacet,
	op crud.Operation,
	item interface{},
) error {
	crudHandler := NewExportsCrud()
	switch op {
	case crud.Create:
		return crudHandler.Create(item)
	case crud.Update:
		return crudHandler.Update(item)
	case crud.Delete:
		return crudHandler.Delete(item)
	case crud.Remove:
		return crudHandler.Delete(item)
	default:
		return fmt.Errorf("operation %s not supported for exports", op)
	}
}
