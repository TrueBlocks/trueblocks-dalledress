package names

import (
	"fmt"
	"sync"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/base"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/output"
	coreTypes "github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/types"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/logging"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/preferences"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/store"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/types"
	sdk "github.com/TrueBlocks/trueblocks-sdk/v5"
)

type Name = coreTypes.Name

var (
	namesStore *store.Store[Name]
	namesMu    sync.Mutex
)

func GetNamesStore() *store.Store[Name] {
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
				// Create structured error with proper context
				wrappedErr := types.NewSDKError("names", types.DataFacet("NamesAll"), "fetch", err)
				logging.LogBackend(fmt.Sprintf("Names SDK query error: %v", wrappedErr))
				return wrappedErr
			}
			logging.LogBackend("The names query function returned without an error.")
			return nil
		}

		processFunc := func(itemIntf interface{}) *Name {
			if name, ok := itemIntf.(*Name); ok {
				return name
			}
			return nil
		}

		mappingFunc := func(item *Name) (key interface{}, includeInMap bool) {
			if item == nil || item.Address.IsZero() {
				return nil, false
			}
			return item.Address, true
		}

		namesStore = store.NewStore(GetStoreName(NamesAll), queryFunc, processFunc, mappingFunc)
	}

	return namesStore
}

func GetStoreName(dataFacet types.DataFacet) string {
	switch dataFacet {
	case NamesAll, NamesCustom, NamesPrefund, NamesRegular, NamesBaddress:
		return "names-list"
	default:
		return ""
	}
}

func (nc *NamesCollection) NameFromAddress(address base.Address) (*Name, bool) {
	return namesStore.GetItemFromMap(address)
}
