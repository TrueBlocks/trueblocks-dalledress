// NAMES_ROUTE
package names

import (
	"fmt"
	"sync"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/base"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/logger"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/output"
	coreTypes "github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/types"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/preferences"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/store"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/types"
	sdk "github.com/TrueBlocks/trueblocks-sdk/v5"
)

var (
	namesStore *store.Store[coreTypes.Name]
	namesMu    sync.Mutex
)

func GetNamesStore() *store.Store[coreTypes.Name] {
	namesMu.Lock()
	defer namesMu.Unlock()

	if namesStore == nil {
		queryFunc := func(ctx *output.RenderCtx) error {
			chainName := preferences.GetChain()
			listOpts := sdk.NamesOptions{
				Globals:   sdk.Globals{Verbose: true, Chain: chainName},
				RenderCtx: ctx,
				All:       true,
			}
			if _, _, err := listOpts.Names(); err != nil {
				logger.Error(fmt.Sprintf("Shared Names source query error: %v", err))
				return err
			}
			return nil
		}

		processFunc := func(itemIntf interface{}) *coreTypes.Name {
			if name, ok := itemIntf.(*coreTypes.Name); ok {
				return name
			}
			return nil
		}

		storeName := GetStoreName(NamesAll)
		namesStore = store.NewStore(storeName, queryFunc, processFunc, nil)
	}

	return namesStore
}

func GetStoreName(listKind types.ListKind) string {
	switch listKind {
	case NamesAll, NamesCustom, NamesPrefund, NamesRegular, NamesBaddress:
		return "names-list"
	default:
		return ""
	}
}

func (nc *NamesCollection) NameFromAddress(address base.Address) (*coreTypes.Name, error) {
	return namesStore.GetItemFromMap(address.Hex()), nil
}

// NAMES_ROUTE
