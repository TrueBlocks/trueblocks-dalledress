// ADD_ROUTE
package abis

import (
	"fmt"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/logger"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/output"
	coreTypes "github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/types"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/repository"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/types"
	sdk "github.com/TrueBlocks/trueblocks-sdk/v5"
)

// ListKind constants for ABI types
const (
	AbisDownloaded types.ListKind = "Downloaded"
	AbisKnown      types.ListKind = "Known"
	AbisFunctions  types.ListKind = "Functions"
	AbisEvents     types.ListKind = "Events"
)

func init() {
	types.RegisterKind(AbisDownloaded)
	types.RegisterKind(AbisKnown)
	types.RegisterKind(AbisFunctions)
	types.RegisterKind(AbisEvents)
}

// AbisCollection orchestrates all ABI repositories
// (Downloaded, Known, Functions, Events)
type AbisCollection struct {
	downloadedRepo repository.Repository[coreTypes.Abi]
	knownRepo      repository.Repository[coreTypes.Abi]
	functionsRepo  repository.Repository[coreTypes.Function]
	eventsRepo     repository.Repository[coreTypes.Function]
}

func NewAbisCollection() AbisCollection {
	// Downloaded: not known
	downloadedRepo := NewAbisRepository(AbisDownloaded, func(abi *coreTypes.Abi) bool {
		return !abi.IsKnown
	})
	// Known: is known
	knownRepo := NewAbisRepository(AbisKnown, func(abi *coreTypes.Abi) bool {
		return abi.IsKnown
	})
	// Functions: not event
	functionsRepo := NewFunctionsRepository(AbisFunctions, func(item *coreTypes.Function) bool {
		return item.FunctionType != "event"
	})
	// Events: is event
	eventsRepo := NewFunctionsRepository(AbisEvents, func(item *coreTypes.Function) bool {
		return item.FunctionType == "event"
	})

	return AbisCollection{
		downloadedRepo: downloadedRepo,
		knownRepo:      knownRepo,
		functionsRepo:  functionsRepo,
		eventsRepo:     eventsRepo,
	}
}

// AbisRepository wraps BaseRepository for coreTypes.Abi
type AbisRepository struct {
	*repository.BaseRepository[coreTypes.Abi]
}

func NewAbisRepository(
	listKind types.ListKind,
	filterFunc repository.FilterFunc[coreTypes.Abi],
) *AbisRepository {
	processFunc := func(itemIntf interface{}) *coreTypes.Abi {
		itemPtr, ok := itemIntf.(*coreTypes.Abi)
		if !ok {
			return nil
		}
		return itemPtr
	}
	queryFunc := func(renderCtx *output.RenderCtx) {
		listOpts := sdk.AbisOptions{
			Globals:   sdk.Globals{Cache: true, Verbose: true},
			RenderCtx: renderCtx,
		}
		if _, _, err := listOpts.AbisList(); err != nil {
			logger.Error(fmt.Sprintf("AbisRepository query error: %v", err))
		}
	}
	baseRepo := repository.NewBaseRepository(
		listKind,
		filterFunc,
		processFunc,
		queryFunc,
		nil, // no dedup needed
	)
	return &AbisRepository{BaseRepository: baseRepo}
}

// FunctionsRepository wraps BaseRepository for coreTypes.Function
type FunctionsRepository struct {
	*repository.BaseRepository[coreTypes.Function]
}

func NewFunctionsRepository(
	listKind types.ListKind,
	filterFunc repository.FilterFunc[coreTypes.Function],
) *FunctionsRepository {
	processFunc := func(itemIntf interface{}) *coreTypes.Function {
		itemPtr, ok := itemIntf.(*coreTypes.Function)
		if !ok {
			return nil
		}
		return itemPtr
	}
	queryFunc := func(renderCtx *output.RenderCtx) {
		detailOpts := sdk.AbisOptions{
			Globals:   sdk.Globals{Cache: true},
			RenderCtx: renderCtx,
		}
		if _, _, err := detailOpts.AbisDetails(); err != nil {
			logger.Error(fmt.Sprintf("FunctionsRepository query error: %v", err))
		}
	}
	dedupeFunc := func(existing []coreTypes.Function, newItem *coreTypes.Function) bool {
		if newItem == nil {
			return false
		}
		for _, item := range existing {
			if item.Encoding == newItem.Encoding {
				return true
			}
		}
		return false
	}
	baseRepo := repository.NewBaseRepository(
		listKind,
		filterFunc,
		processFunc,
		queryFunc,
		dedupeFunc,
	)
	return &FunctionsRepository{BaseRepository: baseRepo}
}

// ADD_ROUTE
