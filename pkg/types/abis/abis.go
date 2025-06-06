// ADD_ROUTE
package abis

import (
	"fmt"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/logger"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/output"
	coreTypes "github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/types"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/facets"
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

// AbisCollection orchestrates all ABI facets
// (Downloaded, Known, Functions, Events)
type AbisCollection struct {
	downloadedFacet facets.Facet[coreTypes.Abi]
	knownFacet      facets.Facet[coreTypes.Abi]
	functionsFacet  facets.Facet[coreTypes.Function]
	eventsFacet     facets.Facet[coreTypes.Function]
}

func NewAbisCollection() AbisCollection {
	downloadedFacet := NewAbisFacet(AbisDownloaded, func(abi *coreTypes.Abi) bool {
		return !abi.IsKnown
	})
	knownFacet := NewAbisFacet(AbisKnown, func(abi *coreTypes.Abi) bool {
		return abi.IsKnown
	})
	functionsFacet := NewFunctionsFacet(AbisFunctions, func(item *coreTypes.Function) bool {
		return item.FunctionType != "event"
	})
	eventsFacet := NewFunctionsFacet(AbisEvents, func(item *coreTypes.Function) bool {
		return item.FunctionType == "event"
	})

	return AbisCollection{
		downloadedFacet: downloadedFacet,
		knownFacet:      knownFacet,
		functionsFacet:  functionsFacet,
		eventsFacet:     eventsFacet,
	}
}

// AbisRepository wraps BaseRepository for coreTypes.Abi
type AbisFacet struct {
	*facets.BaseFacet[coreTypes.Abi]
}

func NewAbisFacet(
	listKind types.ListKind,
	filterFunc facets.FilterFunc[coreTypes.Abi],
) *AbisFacet {
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
			logger.Error(fmt.Sprintf("AbisFacet query error: %v", err))
		}
	}
	baseRepo := facets.NewBaseFacet(
		listKind,
		filterFunc,
		processFunc,
		queryFunc,
		nil,
	)
	return &AbisFacet{BaseFacet: baseRepo}
}

// FunctionsRepository wraps BaseRepository for coreTypes.Function
type FunctionsFacet struct {
	*facets.BaseFacet[coreTypes.Function]
}

func NewFunctionsFacet(
	listKind types.ListKind,
	filterFunc facets.FilterFunc[coreTypes.Function],
) *FunctionsFacet {
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
			logger.Error(fmt.Sprintf("FunctionsFacet query error: %v", err))
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
	baseRepo := facets.NewBaseFacet(
		listKind,
		filterFunc,
		processFunc,
		queryFunc,
		dedupeFunc,
	)
	return &FunctionsFacet{BaseFacet: baseRepo}
}

// ADD_ROUTE
