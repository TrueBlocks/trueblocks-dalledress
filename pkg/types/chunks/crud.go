package chunks

import (
	"fmt"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/crud"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/types"
)

type ChunksCrud struct{}

func NewChunksCrud() *ChunksCrud {
	return &ChunksCrud{}
}

func (c *ChunksCollection) Crud(
	dataFacet types.DataFacet,
	op crud.Operation,
	item interface{},
) error {
	// crudHandler := NewChunksCrud()
	return fmt.Errorf("operation %s not supported for chunks", op)
}
