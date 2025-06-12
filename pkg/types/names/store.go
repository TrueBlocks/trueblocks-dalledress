// NAMES_ROUTE
package names

import (
	"fmt"
	"sync"

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
	storeMu    sync.Mutex
)

func GetNamesStore() *store.Store[coreTypes.Name] {
	storeMu.Lock()
	defer storeMu.Unlock()

	if namesStore == nil {
		namesStore = store.NewStore(
			GetStoreName(NamesAll),
			func(ctx *output.RenderCtx) error {
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
			},
			func(itemIntf interface{}) *coreTypes.Name {
				if name, ok := itemIntf.(*coreTypes.Name); ok {
					return name
				}
				return nil
			},
		)
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

// NAMES_ROUTE
