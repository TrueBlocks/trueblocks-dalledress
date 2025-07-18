package monitors

import (
	"fmt"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/base"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/crud"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/logging"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/msgs"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/preferences"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/types"
	sdk "github.com/TrueBlocks/trueblocks-sdk/v5"
)

func (c *MonitorsCollection) Crud(
	payload *types.Payload,
	op crud.Operation,
	item interface{},
) error {
	var monitor = &Monitor{Address: base.HexToAddress(payload.Address)}
	if cast, ok := item.(*Monitor); ok && cast != nil {
		monitor = cast
	}

	chain := preferences.GetLastChain()

	var err error
	switch op {
	case crud.Remove:
		opts := sdk.MonitorsOptions{
			Addrs:   []string{monitor.Address.Hex()},
			Remove:  true,
			Globals: sdk.Globals{Chain: chain},
		}
		_, _, err = opts.Monitors()
	case crud.Delete:
		opts := sdk.MonitorsOptions{
			Addrs:   []string{monitor.Address.Hex()},
			Delete:  true,
			Globals: sdk.Globals{Cache: true, Chain: chain},
		}
		_, _, err = opts.Monitors()
	case crud.Undelete:
		opts := sdk.MonitorsOptions{
			Addrs:    []string{monitor.Address.Hex()},
			Undelete: true,
			Globals:  sdk.Globals{Cache: true, Chain: chain},
		}
		_, _, err = opts.Monitors()
	default:
		logging.LogBackend(fmt.Sprintf("Monitor operation %s not implemented for address: %s", op, monitor.Address))
		return fmt.Errorf("operation %s not yet implemented for Monitors", op)
	}

	if err != nil {
		return err
	}

	store := c.monitorsFacet.GetStore()
	store.UpdateData(func(data []*Monitor) []*Monitor {
		return c.updateMonitorInData(data, monitor, op)
	})
	c.monitorsFacet.SyncWithStore()

	switch op {
	case crud.Remove:
		msgs.EmitStatus(fmt.Sprintf("removed monitor for address: %s", monitor.Address))
		logging.LogBackend(fmt.Sprintf("Removed monitor for address: %s", monitor.Address))
	case crud.Delete:
		msgs.EmitStatus(fmt.Sprintf("deleted monitor for address: %s", monitor.Address))
		logging.LogBackend(fmt.Sprintf("Deleted monitor for address: %s", monitor.Address))
	case crud.Undelete:
		msgs.EmitStatus(fmt.Sprintf("undeleted monitor for address: %s", monitor.Address))
		logging.LogBackend(fmt.Sprintf("Undeleted monitor for address: %s", monitor.Address))
	}

	return nil
}

func (c *MonitorsCollection) updateMonitorInData(data []*Monitor, monitor *Monitor, op crud.Operation) []*Monitor {
	switch op {
	case crud.Remove:
		result := make([]*Monitor, 0, len(data))
		for _, m := range data {
			if m.Address != monitor.Address {
				result = append(result, m)
			}
		}
		return result
	case crud.Delete:
		for _, m := range data {
			if m.Address == monitor.Address {
				m.Deleted = true
				break
			}
		}
		return data
	case crud.Undelete:
		for _, m := range data {
			if m.Address == monitor.Address {
				m.Deleted = false
				break
			}
		}
		return data
	default:
		logging.LogBackend(fmt.Sprintf("Monitor operation %s not implemented for address: %s", op, monitor.Address))
		return data
	}
}

func (c *MonitorsCollection) Clean(addresses []string) error {
	chain := preferences.GetLastChain()

	opts := sdk.MonitorsOptions{
		Globals: sdk.Globals{Cache: true, Chain: chain},
	}

	if len(addresses) > 0 {
		opts.Addrs = addresses
	}

	cleanResult, _, err := opts.MonitorsClean()
	if err != nil {
		return err
	}

	if len(addresses) > 0 {
		msgs.EmitStatus(fmt.Sprintf("cleaned %d monitor(s)", len(addresses)))
		logging.LogBackend(fmt.Sprintf("Cleaned monitors for addresses: %v", addresses))
	} else {
		msgs.EmitStatus(fmt.Sprintf("cleaned all monitors, processed %d items", len(cleanResult)))
		logging.LogBackend("Cleaned all monitors")
	}

	c.LoadData(MonitorsMonitors)
	return nil
}
