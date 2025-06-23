package names

import (
	"sync/atomic"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/base"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/crud"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/msgs"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/types"
	sdk "github.com/TrueBlocks/trueblocks-sdk/v5"
)

var namesLock atomic.Int32

func (c *NamesCollection) Crud(
	payload types.Payload,
	op crud.Operation,
	item interface{},
) error {
	dataFacet := payload.DataFacet

	var name = &Name{Address: base.HexToAddress(payload.Address)}
	if cast, ok := item.(*Name); ok && cast != nil {
		name = cast
	}

	if !namesLock.CompareAndSwap(0, 1) {
		return nil
	}
	defer namesLock.Store(0)

	name.IsCustom = true

	cd := crud.CrudFromName(*name)
	opts := sdk.NamesOptions{
		Globals: sdk.Globals{
			Chain: "mainnet",
		},
	}

	if _, _, err := opts.ModifyName(op, cd); err != nil {
		msgs.EmitError("Crud", err)
		return err
	}

	// TODO: See the AbisCollection for in-memory cache updating code instead of full Reset.
	c.Reset(dataFacet)
	return nil
}

// TODO: Consider adding batch operations for Names, similar to MonitorsCollection.Clean (e.g., batch delete).
