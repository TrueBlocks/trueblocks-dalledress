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
	payload types.Payload,
	op crud.Operation,
	item interface{},
) error {
	return fmt.Errorf("operation %s not supported for exports", op)
}
