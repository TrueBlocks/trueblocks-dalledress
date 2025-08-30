package dalledress

import (
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/crud"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/types"
)

// Crud implements CRUD operations for chunks - placeholder implementation
func (c *DalleDressCollection) Crud(
	payload *types.Payload,
	op crud.Operation,
	item interface{},
) error {
	switch v := item.(type) {
	case *Series:
		return c.seriesCrud(payload, op, v)
	// Add other facet types here as needed
	default:
		// Placeholder for other facets
		return nil
	}
}
