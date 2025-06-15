package monitors

import (
	"fmt"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/crud"
	coreTypes "github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/types"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/logging"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/msgs"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/preferences"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/types"
	sdk "github.com/TrueBlocks/trueblocks-sdk/v5"
)

func (c *MonitorsCollection) Crud(
	listKind types.ListKind,
	op crud.Operation,
	monitor *coreTypes.Monitor,
) error {
	chainName := preferences.GetChain()

	switch op {
	case crud.Remove:
		opts := sdk.MonitorsOptions{
			Addrs:   []string{monitor.Address.Hex()},
			Remove:  true,
			Globals: sdk.Globals{Chain: chainName},
		}
		if _, _, err := opts.Monitors(); err != nil {
			return err
		}

		msgs.EmitStatus(fmt.Sprintf("removed monitor for address: %s", monitor.Address))
		logging.LogBackend(fmt.Sprintf("Removed monitor for address: %s", monitor.Address))

		// TODO: See the AbisCollection for in-memory cache updating code instead of full Reset.
		c.Reset(MonitorsList)
		return nil

	case crud.Delete:
		opts := sdk.MonitorsOptions{
			Addrs:   []string{monitor.Address.Hex()},
			Delete:  true,
			Globals: sdk.Globals{Cache: true, Chain: chainName},
		}
		if _, _, err := opts.Monitors(); err != nil {
			return err
		}

		msgs.EmitStatus(fmt.Sprintf("deleted monitor for address: %s", monitor.Address))
		logging.LogBackend(fmt.Sprintf("Deleted monitor for address: %s", monitor.Address))

		// TODO: See the AbisCollection for in-memory cache updating code instead of full Reset.
		c.Reset(MonitorsList)
		return nil

	case crud.Undelete:
		opts := sdk.MonitorsOptions{
			Addrs:    []string{monitor.Address.Hex()},
			Undelete: true,
			Globals:  sdk.Globals{Cache: true, Chain: chainName},
		}
		if _, _, err := opts.Monitors(); err != nil {
			return err
		}

		msgs.EmitStatus(fmt.Sprintf("undeleted monitor for address: %s", monitor.Address))
		logging.LogBackend(fmt.Sprintf("Undeleted monitor for address: %s", monitor.Address))

		// TODO: See the AbisCollection for in-memory cache updating code instead of full Reset.
		c.Reset(MonitorsList)
		return nil

	default:
		logging.LogBackend(fmt.Sprintf("Monitor operation %s not implemented for address: %s", op, monitor.Address))
		return fmt.Errorf("operation %s not yet implemented for Monitors", op)
	}
}

func (c *MonitorsCollection) Clean(addresses []string) error {
	chainName := preferences.GetChain()

	opts := sdk.MonitorsOptions{
		Globals: sdk.Globals{Cache: true, Chain: chainName},
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

	// TODO: See the AbisCollection for in-memory cache updating code instead of full Reset.
	c.Reset(MonitorsList)
	return nil
}
