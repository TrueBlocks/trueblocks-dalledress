package names

import (
	"fmt"
	"sync/atomic"

	"github.com/TrueBlocks/trueblocks-dalledress/pkg/logging"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/msgs"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/types"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/base"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/crud"
	sdk "github.com/TrueBlocks/trueblocks-sdk/v5"
)

var namesLock atomic.Int32

func (c *NamesCollection) Crud(
	payload *types.Payload,
	op crud.Operation,
	item *Name,
) error {
	dataFacet := payload.DataFacet
	logging.LogBackend(fmt.Sprintf("[NAMES CRUD] Enter Crud op=%s facet=%s payloadAddr=%s chain=%s itemNil=%t", op, dataFacet, payload.Address, payload.Chain, item == nil))

	var name = &Name{Address: base.HexToAddress(payload.Address)}
	if item != nil {
		name = item
	}
	logging.LogBackend(fmt.Sprintf("[NAMES CRUD] Prepared name address=%s incomingOp=%s isCustomInitial=%t deletedInitial=%t", name.Address, op, name.Parts != 0, name.Deleted))

	if !namesLock.CompareAndSwap(0, 1) {
		logging.LogBackend("[NAMES CRUD] Lock busy, dropping duplicate request")
		return nil
	}
	defer func() {
		namesLock.Store(0)
		logging.LogBackend("[NAMES CRUD] Released lock")
	}()

	name.IsCustom = true
	logging.LogBackend(fmt.Sprintf("[NAMES CRUD] Marked name as custom address=%s", name.Address))

	cd := crud.CrudFromName(*name)
	opts := sdk.NamesOptions{
		Globals: sdk.Globals{
			Chain: "mainnet",
		},
	}
	logging.LogBackend(fmt.Sprintf("[NAMES CRUD] Calling ModifyName op=%s address=%s", op, name.Address))

	if _, _, err := opts.ModifyName(op, cd); err != nil {
		logging.LogBackend(fmt.Sprintf("[NAMES CRUD] ModifyName error op=%s address=%s err=%v", op, name.Address, err))
		msgs.EmitError("Crud", err)
		return err
	}
	logging.LogBackend(fmt.Sprintf("[NAMES CRUD] ModifyName success op=%s address=%s", op, name.Address))

	if op == crud.Autoname {
		logging.LogBackend(fmt.Sprintf("[NAMES CRUD] Autoname triggered full facet reset facet=%s", dataFacet))
		c.Reset(dataFacet)
		msgs.EmitStatus(fmt.Sprintf("completed %s operation for name: %s", op, name.Address))
		return nil
	}

	switch dataFacet {
	case NamesAll:
		store := c.allFacet.GetStore()
		before := 0
		after := 0
		store.UpdateData(func(data []*Name) []*Name {
			before = len(data)
			updated := c.updateNameInData(data, name, op)
			after = len(updated)
			return updated
		})
		c.allFacet.SyncWithStore()
		logging.LogBackend(fmt.Sprintf("[NAMES CRUD] Synced facet=all before=%d after=%d", before, after))
	case NamesCustom:
		store := c.customFacet.GetStore()
		before := 0
		after := 0
		store.UpdateData(func(data []*Name) []*Name {
			before = len(data)
			updated := c.updateNameInData(data, name, op)
			after = len(updated)
			return updated
		})
		c.customFacet.SyncWithStore()
		logging.LogBackend(fmt.Sprintf("[NAMES CRUD] Synced facet=custom before=%d after=%d", before, after))
	case NamesPrefund:
		store := c.prefundFacet.GetStore()
		before := 0
		after := 0
		store.UpdateData(func(data []*Name) []*Name {
			before = len(data)
			updated := c.updateNameInData(data, name, op)
			after = len(updated)
			return updated
		})
		c.prefundFacet.SyncWithStore()
		logging.LogBackend(fmt.Sprintf("[NAMES CRUD] Synced facet=prefund before=%d after=%d", before, after))
	case NamesRegular:
		store := c.regularFacet.GetStore()
		before := 0
		after := 0
		store.UpdateData(func(data []*Name) []*Name {
			before = len(data)
			updated := c.updateNameInData(data, name, op)
			after = len(updated)
			return updated
		})
		c.regularFacet.SyncWithStore()
		logging.LogBackend(fmt.Sprintf("[NAMES CRUD] Synced facet=regular before=%d after=%d", before, after))
	case NamesBaddress:
		store := c.baddressFacet.GetStore()
		before := 0
		after := 0
		store.UpdateData(func(data []*Name) []*Name {
			before = len(data)
			updated := c.updateNameInData(data, name, op)
			after = len(updated)
			return updated
		})
		c.baddressFacet.SyncWithStore()
		logging.LogBackend(fmt.Sprintf("[NAMES CRUD] Synced facet=baddress before=%d after=%d", before, after))
	default:
		logging.LogBackend(fmt.Sprintf("[NAMES CRUD] Unknown facet=%s no in-memory update performed", dataFacet))
	}

	msgs.EmitStatus(fmt.Sprintf("completed %s operation for name: %s", op, name.Address))
	logging.LogBackend(fmt.Sprintf("[NAMES CRUD] Completed op=%s facet=%s address=%s", op, dataFacet, name.Address))
	return nil
}

// updateNameInData handles the in-memory data update logic for all CRUD operations
func (c *NamesCollection) updateNameInData(data []*Name, name *Name, op crud.Operation) []*Name {
	switch op {
	case crud.Remove:
		logging.LogBackend(fmt.Sprintf("[NAMES CRUD] updateNameInData remove address=%s currentLen=%d", name.Address, len(data)))
		result := make([]*Name, 0, len(data))
		for _, n := range data {
			if n.Address != name.Address {
				result = append(result, n)
			}
		}
		logging.LogBackend(fmt.Sprintf("[NAMES CRUD] updateNameInData remove resultLen=%d", len(result)))
		return result
	case crud.Create:
		logging.LogBackend(fmt.Sprintf("[NAMES CRUD] updateNameInData create address=%s", name.Address))
		for _, n := range data {
			if n.Address == name.Address {
				*n = *name
				logging.LogBackend("[NAMES CRUD] updateNameInData create replaced existing")
				return data
			}
		}
		logging.LogBackend("[NAMES CRUD] updateNameInData create appended new")
		return append(data, name)
	case crud.Update:
		logging.LogBackend(fmt.Sprintf("[NAMES CRUD] updateNameInData update address=%s", name.Address))
		for _, n := range data {
			if n.Address == name.Address {
				*n = *name
				logging.LogBackend("[NAMES CRUD] updateNameInData update applied")
				break
			}
		}
		return data
	case crud.Delete:
		logging.LogBackend(fmt.Sprintf("[NAMES CRUD] updateNameInData delete address=%s", name.Address))
		for _, n := range data {
			if n.Address == name.Address {
				n.Deleted = true
				logging.LogBackend("[NAMES CRUD] updateNameInData marked deleted")
				break
			}
		}
		return data
	case crud.Undelete:
		logging.LogBackend(fmt.Sprintf("[NAMES CRUD] updateNameInData undelete address=%s", name.Address))
		for _, n := range data {
			if n.Address == name.Address {
				n.Deleted = false
				logging.LogBackend("[NAMES CRUD] updateNameInData unmarked deleted")
				break
			}
		}
		return data
	default:
		logging.LogBackend(fmt.Sprintf("[NAMES CRUD] updateNameInData no-op op=%s address=%s", op, name.Address))
		return data
	}
}

// TODO: Consider adding batch operations for Names, similar to MonitorsCollection.Clean (e.g., batch delete).
