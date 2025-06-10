package enhancedcollection

import (
	"fmt"
	"strings"

	coreTypes "github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/types"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/enhancedfacet"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/facets"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/logging"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/msgs"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/types"
)

// // Define ListKind constants locally to avoid import cycle
// // These must match the values in pkg/types/abis/abis.go
// const (
// 	AbisDownloaded types.ListKind = "Downloaded"
// 	AbisKnown      types.ListKind = "Known"
// 	AbisFunctions  types.ListKind = "Functions"
// 	AbisEvents     types.ListKind = "Events"
// )

const (
	AbisDownloaded types.ListKind = "Downloaded"
	AbisKnown      types.ListKind = "Known"
	AbisFunctions  types.ListKind = "Functions"
	AbisEvents     types.ListKind = "Events"
)

// Register list kinds
func init() {
	types.RegisterKind(AbisDownloaded)
	types.RegisterKind(AbisKnown)
	types.RegisterKind(AbisFunctions)
	types.RegisterKind(AbisEvents)
}

// AbisPage represents a paginated view of ABIs or functions
// for a given ListKind (Downloaded, Known, Functions, Events)
// This mirrors the original AbisPage struct for API compatibility
type AbisPage struct {
	Kind          types.ListKind       `json:"kind"`
	Abis          []coreTypes.Abi      `json:"abis,omitempty"`
	Functions     []coreTypes.Function `json:"functions,omitempty"`
	TotalItems    int                  `json:"totalItems"`
	ExpectedTotal int                  `json:"expectedTotal"`
	IsFetching    bool                 `json:"isFetching"`
	State         facets.LoadState     `json:"state"`
}

// mapLoadState converts enhancedfacet.LoadState to facets.LoadState for compatibility
func mapLoadState(enhancedState enhancedfacet.LoadState) facets.LoadState {
	switch enhancedState {
	case enhancedfacet.StateStale:
		return facets.StateStale
	case enhancedfacet.StateFetching:
		return facets.StateFetching
	case enhancedfacet.StateLoaded:
		return facets.StateLoaded
	case enhancedfacet.StateError:
		return facets.StateError
	case enhancedfacet.StatePartial:
		return facets.StatePartial
	default:
		return facets.StateStale
	}
}

// EnhancedAbisCollection provides access to ABIs data using the enhanced architecture
type EnhancedAbisCollection struct {
	downloadedFacet *enhancedfacet.BaseFacet[coreTypes.Abi]
	knownFacet      *enhancedfacet.BaseFacet[coreTypes.Abi]
	functionsFacet  *enhancedfacet.BaseFacet[coreTypes.Function]
	eventsFacet     *enhancedfacet.BaseFacet[coreTypes.Function]
}

// NewEnhancedAbisCollection creates a new Abis collection using the enhanced architecture
func NewEnhancedAbisCollection() *EnhancedAbisCollection {
	// Get the shared stores
	abisListStore := GetEnhancedListStore()
	abisDetailStore := GetEnhancedDetailStore()

	// Create Downloaded facet
	downloadedFacet := enhancedfacet.NewBaseFacet(
		AbisDownloaded,
		func(abi *coreTypes.Abi) bool { return !abi.IsKnown },
		nil,
		abisListStore,
	)

	// Create Known facet
	knownFacet := enhancedfacet.NewBaseFacet(
		AbisKnown,
		func(abi *coreTypes.Abi) bool { return abi.IsKnown },
		nil,
		abisListStore,
	)

	functionsFacet := enhancedfacet.NewBaseFacet(
		AbisFunctions,
		func(fn *coreTypes.Function) bool { return fn.FunctionType != "event" },
		func(existing []*coreTypes.Function, newItem *coreTypes.Function) bool {
			for _, existingFn := range existing {
				if existingFn.Encoding == newItem.Encoding {
					return true
				}
			}
			return false
		},
		abisDetailStore,
	)

	eventsFacet := enhancedfacet.NewBaseFacet(
		AbisEvents,
		func(fn *coreTypes.Function) bool { return fn.FunctionType == "event" },
		func(existing []*coreTypes.Function, newItem *coreTypes.Function) bool {
			for _, existingFn := range existing {
				if existingFn.Encoding == newItem.Encoding {
					return true // It's a duplicate
				}
			}
			return false
		},
		abisDetailStore,
	)

	return &EnhancedAbisCollection{
		downloadedFacet: downloadedFacet,
		knownFacet:      knownFacet,
		functionsFacet:  functionsFacet,
		eventsFacet:     eventsFacet,
	}
}

// LoadData loads data for a specific facet
func (c *EnhancedAbisCollection) LoadData(listKind types.ListKind) {
	if !c.NeedsUpdate(listKind) {
		return
	}

	switch listKind {
	case AbisDownloaded: // Use local constant
		go func() {
			if result, err := c.downloadedFacet.Load(); err != nil {
				logging.LogError("LoadData.AbisDownloaded from store: %v", err, enhancedfacet.ErrAlreadyLoading)
			} else {
				logging.LogBackend(fmt.Sprintf("LoadData: kind: %s currentCount: %d expectedTotal: %d", result.Payload.ListKind, result.Payload.CurrentCount, result.Payload.ExpectedTotal))
				msgs.EmitLoaded("downloaded", result.Payload)
			}
		}()
	case AbisKnown: // Use local constant
		go func() {
			if result, err := c.knownFacet.Load(); err != nil {
				logging.LogError("LoadData.AbisKnown from store: %v", err, enhancedfacet.ErrAlreadyLoading)
			} else {
				logging.LogBackend(fmt.Sprintf("LoadData: kind: %s currentCount: %d expectedTotal: %d", result.Payload.ListKind, result.Payload.CurrentCount, result.Payload.ExpectedTotal))
				msgs.EmitLoaded("known", result.Payload)
			}
		}()
	case AbisFunctions: // Use local constant
		go func() {
			if result, err := c.functionsFacet.Load(); err != nil {
				logging.LogError("LoadData.AbisFunctions from store: %v", err, enhancedfacet.ErrAlreadyLoading)
			} else {
				logging.LogBackend(fmt.Sprintf("LoadData: kind: %s currentCount: %d expectedTotal: %d", result.Payload.ListKind, result.Payload.CurrentCount, result.Payload.ExpectedTotal))
				msgs.EmitLoaded("functions", result.Payload)
			}
		}()
	case AbisEvents: // Use local constant
		go func() {
			if result, err := c.eventsFacet.Load(); err != nil {
				logging.LogError("LoadData.AbisEvents from store: %v", err, enhancedfacet.ErrAlreadyLoading)
			} else {
				logging.LogBackend(fmt.Sprintf("LoadData: kind: %s currentCount: %d expectedTotal: %d", result.Payload.ListKind, result.Payload.CurrentCount, result.Payload.ExpectedTotal))
				msgs.EmitLoaded("events", result.Payload)
			}
		}()
	default:
		logging.LogError("LoadData: unexpected list kind: %v", fmt.Errorf("invalid list kind: %s", listKind), nil)
	}
}

// Reset resets a specific facet
func (c *EnhancedAbisCollection) Reset(listKind types.ListKind) {
	switch listKind {
	case AbisDownloaded: // Use local constant
		c.downloadedFacet.GetStore().Reset()
	case AbisKnown: // Use local constant
		c.knownFacet.GetStore().Reset()
	case AbisFunctions: // Use local constant
		c.functionsFacet.GetStore().Reset()
	case AbisEvents: // Use local constant
		c.eventsFacet.GetStore().Reset()
	}
}

// NeedsUpdate checks if a specific facet needs updating
func (c *EnhancedAbisCollection) NeedsUpdate(listKind types.ListKind) bool {
	switch listKind {
	case AbisDownloaded: // Use local constant
		return c.downloadedFacet.NeedsUpdate()
	case AbisKnown: // Use local constant
		return c.knownFacet.NeedsUpdate()
	case AbisFunctions: // Use local constant
		return c.functionsFacet.NeedsUpdate()
	case AbisEvents: // Use local constant
		return c.eventsFacet.NeedsUpdate()
	}
	return false
}

// GetDownloadedPage returns a page of downloaded ABIs
func (c *EnhancedAbisCollection) GetDownloadedPage(first, pageSize int) ([]coreTypes.Abi, int, enhancedfacet.LoadState) {
	result, err := c.downloadedFacet.GetPage(first, pageSize, nil, nil, nil)
	if err != nil {
		return nil, 0, enhancedfacet.StateError
	}
	return result.Items, result.TotalItems, result.State
}

// GetKnownPage returns a page of known ABIs
func (c *EnhancedAbisCollection) GetKnownPage(first, pageSize int) ([]coreTypes.Abi, int, enhancedfacet.LoadState) {
	result, err := c.knownFacet.GetPage(first, pageSize, nil, nil, nil)
	if err != nil {
		return nil, 0, enhancedfacet.StateError
	}
	return result.Items, result.TotalItems, result.State
}

// GetFunctionsPage returns a page of functions
func (c *EnhancedAbisCollection) GetFunctionsPage(first, pageSize int) ([]coreTypes.Function, int, enhancedfacet.LoadState) {
	result, err := c.functionsFacet.GetPage(first, pageSize, nil, nil, nil)
	if err != nil {
		return nil, 0, enhancedfacet.StateError
	}
	return result.Items, result.TotalItems, result.State
}

// GetEventsPage returns a page of events
func (c *EnhancedAbisCollection) GetEventsPage(first, pageSize int) ([]coreTypes.Function, int, enhancedfacet.LoadState) {
	result, err := c.eventsFacet.GetPage(first, pageSize, nil, nil, nil)
	if err != nil {
		return nil, 0, enhancedfacet.StateError
	}
	return result.Items, result.TotalItems, result.State
}

// GetDownloadedState returns the state of the downloaded facet
func (c *EnhancedAbisCollection) GetDownloadedState() enhancedfacet.LoadState {
	return c.downloadedFacet.GetState()
}

// GetKnownState returns the state of the known facet
func (c *EnhancedAbisCollection) GetKnownState() enhancedfacet.LoadState {
	return c.knownFacet.GetState()
}

// GetFunctionsState returns the state of the functions facet
func (c *EnhancedAbisCollection) GetFunctionsState() enhancedfacet.LoadState {
	return c.functionsFacet.GetState()
}

// GetEventsState returns the state of the events facet
func (c *EnhancedAbisCollection) GetEventsState() enhancedfacet.LoadState {
	return c.eventsFacet.GetState()
}

// GetAllAbis returns all ABIs from the downloaded and known facets
func (c *EnhancedAbisCollection) GetAllAbis() ([]coreTypes.Abi, error) {
	var allAbis []coreTypes.Abi

	// Get downloaded ABIs
	downloadedAbis, _, _ := c.GetDownloadedPage(0, 1000000) // Large page size to fetch all
	allAbis = append(allAbis, downloadedAbis...)

	// Get known ABIs
	knownAbis, _, _ := c.GetKnownPage(0, 1000000) // Large page size to fetch all
	allAbis = append(allAbis, knownAbis...)

	return allAbis, nil
}

// GetAllFunctions returns all Functions from the functions facet
func (c *EnhancedAbisCollection) GetAllFunctions() ([]coreTypes.Function, error) {
	functions, _, _ := c.GetFunctionsPage(0, 1000000) // Large page size to fetch all
	return functions, nil
}

// GetAllEvents returns all Events from the events facet
func (c *EnhancedAbisCollection) GetAllEvents() ([]coreTypes.Function, error) {
	events, _, _ := c.GetEventsPage(0, 1000000) // Large page size to fetch all
	return events, nil
}

// FilterAbisByName filters ABIs by name from the downloaded and known facets
func (c *EnhancedAbisCollection) FilterAbisByName(name string) ([]coreTypes.Abi, error) {
	var filteredAbis []coreTypes.Abi

	// Get all ABIs
	allAbis, err := c.GetAllAbis()
	if err != nil {
		return nil, err
	}

	// Filter by name
	for _, abi := range allAbis {
		if strings.Contains(abi.Name, name) {
			filteredAbis = append(filteredAbis, abi)
		}
	}

	return filteredAbis, nil
}

// FilterFunctionsByName filters Functions by name from the functions facet
func (c *EnhancedAbisCollection) FilterFunctionsByName(name string) ([]coreTypes.Function, error) {
	var filteredFunctions []coreTypes.Function

	// Get all Functions
	allFunctions, err := c.GetAllFunctions()
	if err != nil {
		return nil, err
	}

	// Filter by name
	for _, function := range allFunctions {
		if strings.Contains(function.Name, name) {
			filteredFunctions = append(filteredFunctions, function)
		}
	}

	return filteredFunctions, nil
}

// FilterEventsByName filters Events by name from the events facet
func (c *EnhancedAbisCollection) FilterEventsByName(name string) ([]coreTypes.Function, error) {
	var filteredEvents []coreTypes.Function

	// Get all Events
	allEvents, err := c.GetAllEvents()
	if err != nil {
		return nil, err
	}

	// Filter by name
	for _, event := range allEvents {
		if strings.Contains(event.Name, name) {
			filteredEvents = append(filteredEvents, event)
		}
	}

	return filteredEvents, nil
}

// GetAbiByAddress retrieves an ABI by its address from the downloaded or known facets
func (c *EnhancedAbisCollection) GetAbiByAddress(address string) (*coreTypes.Abi, error) {
	// Check in downloaded ABIs
	downloadedAbis, _, _ := c.GetDownloadedPage(0, 1000000)
	for _, abi := range downloadedAbis {
		if abi.Address.Hex() == address {
			return &abi, nil
		}
	}

	// Check in known ABIs
	knownAbis, _, _ := c.GetKnownPage(0, 1000000)
	for _, abi := range knownAbis {
		if abi.Address.Hex() == address {
			return &abi, nil
		}
	}

	return nil, fmt.Errorf("ABI not found")
}

// GetFunctionBySignature retrieves a Function by its signature from the functions facet
func (c *EnhancedAbisCollection) GetFunctionBySignature(signature string) (*coreTypes.Function, error) {
	functions, _, _ := c.GetFunctionsPage(0, 1000000)
	for _, function := range functions {
		if function.Signature == signature {
			return &function, nil
		}
	}

	return nil, fmt.Errorf("function not found")
}

// GetEventBySignature retrieves an Event by its signature from the events facet
func (c *EnhancedAbisCollection) GetEventBySignature(signature string) (*coreTypes.Function, error) {
	events, _, _ := c.GetEventsPage(0, 1000000)
	for _, event := range events {
		if event.Signature == signature {
			return &event, nil
		}
	}

	return nil, fmt.Errorf("event not found")
}

// GetPage returns a filtered, sorted, and paginated AbisPage for the given ListKind
// This method provides full API compatibility with the original AbisCollection
func (c *EnhancedAbisCollection) GetPage(
	first, pageSize int,
	filterString string,
	sortSpec string,
	listKind types.ListKind,
) (*AbisPage, error) {
	page := &AbisPage{
		Kind: listKind,
	}

	switch listKind {
	case AbisDownloaded: // Use local constant
		items, total, state := c.GetDownloadedPage(first, pageSize)
		page.Abis = items
		page.TotalItems = total
		page.State = mapLoadState(state)
		page.IsFetching = c.downloadedFacet.IsFetching()
		page.ExpectedTotal = c.downloadedFacet.ExpectedCount()

	case AbisKnown: // Use local constant
		items, total, state := c.GetKnownPage(first, pageSize)
		page.Abis = items
		page.TotalItems = total
		page.State = mapLoadState(state)
		page.IsFetching = c.knownFacet.IsFetching()
		page.ExpectedTotal = c.knownFacet.ExpectedCount()

	case AbisFunctions, AbisEvents: // Use local constants
		var items []coreTypes.Function
		var total int
		var state enhancedfacet.LoadState
		var facet *enhancedfacet.BaseFacet[coreTypes.Function]

		if listKind == AbisFunctions { // Use local constant
			items, total, state = c.GetFunctionsPage(first, pageSize)
			facet = c.functionsFacet
		} else { // Must be AbisEvents
			items, total, state = c.GetEventsPage(first, pageSize)
			facet = c.eventsFacet
		}

		page.Functions = items
		page.TotalItems = total
		page.State = mapLoadState(state)
		page.IsFetching = facet.IsFetching()
		page.ExpectedTotal = facet.ExpectedCount()

	default:
		return nil, fmt.Errorf("GetPage: unexpected list kind: %v", listKind)
	}

	// TODO: Implement filtering and sorting if filterString or sortSpec are provided
	// For now, we just log if they are present as a reminder.
	if filterString != "" {
		logging.LogInfo("GetPage: filterString is present but not yet implemented: %s", filterString)
	}
	if sortSpec != "" {
		logging.LogInfo("GetPage: sortSpec is present but not yet implemented: %s", sortSpec)
	}

	return page, nil
}
