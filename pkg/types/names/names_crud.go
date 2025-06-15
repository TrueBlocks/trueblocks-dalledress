package names

import (
	"fmt"
	"sync/atomic"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/crud"
	coreTypes "github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/types"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/msgs"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/types"
	sdk "github.com/TrueBlocks/trueblocks-sdk/v5"
)

var namesLock atomic.Int32

func (nc *NamesCollection) Crud(
	listKind types.ListKind,
	op crud.Operation,
	name *coreTypes.Name,
) error {
	if name == nil || name.Address.IsZero() {
		return fmt.Errorf("Crud operation requires a valid name with a non-zero address")
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
	nc.Reset(listKind)
	return nil
}

// TODO: Consider adding batch operations for Names, similar to MonitorsCollection.Clean (e.g., batch delete).
