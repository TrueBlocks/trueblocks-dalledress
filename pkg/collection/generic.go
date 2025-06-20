package collection

import (
	"fmt"
	"sync"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/crud"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/facets"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/types"
	sdk "github.com/TrueBlocks/trueblocks-sdk/v5"
)

// PageFactory creates a concrete Page from facet PageResult
// This allows each collection type to create its own Page implementation
type PageFactory[T any] interface {
	CreatePage(dataFacet types.DataFacet, pageResult *facets.PageResult[T], facet *facets.Facet[T]) (types.Page, error)
}

// FilterFunc creates a filter function from a string filter
type FilterFunc[T any] interface {
	CreateFilter(filter string) func(*T) bool
}

// SortFunc creates a sort function for the given type
type SortFunc[T any] interface {
	CreateSort() func([]T, sdk.SortSpec) error
}

// CrudHandler handles CRUD operations for a specific type
type CrudHandler[T any] interface {
	Crud(op crud.Operation, item T) error
}

// HandlerInterface provides a type-erased interface for collection operations
// This allows us to store different typed handlers in the same registry
type HandlerInterface interface {
	LoadData() error
	Reset()
	NeedsUpdate() bool
	GetPageAny(first, pageSize int, sort sdk.SortSpec, filter string) (types.Page, error)
	CrudAny(op crud.Operation, item interface{}) error
	GetStoreName() string
	GetFacet() types.DataFacet
}

// TypedHandler provides a type-safe wrapper around facets for a specific type T
// This bridges the gap between the generic interface and strongly-typed facet operations
type TypedHandler[T any] struct {
	dataFacet   types.DataFacet
	storeName   string
	facet       *facets.Facet[T]
	pageFactory PageFactory[T]
	filterFunc  FilterFunc[T]
	sortFunc    SortFunc[T]
	crudHandler CrudHandler[T]
}

// NewTypedHandler creates a new typed handler for a specific type and dataFacet
func NewTypedHandler[T any](
	dataFacet types.DataFacet,
	storeName string,
	facet *facets.Facet[T],
	pageFactory PageFactory[T],
	filterFunc FilterFunc[T],
	sortFunc SortFunc[T],
	crudHandler CrudHandler[T],
) *TypedHandler[T] {
	return &TypedHandler[T]{
		dataFacet:   dataFacet,
		storeName:   storeName,
		facet:       facet,
		pageFactory: pageFactory,
		filterFunc:  filterFunc,
		sortFunc:    sortFunc,
		crudHandler: crudHandler,
	}
}

// Implement HandlerInterface for TypedHandler
func (th *TypedHandler[T]) LoadData() error {
	// Use the facet's Load method which returns a streaming result
	return th.facet.Load()
}

func (th *TypedHandler[T]) Reset() {
	th.facet.Reset()
}

func (th *TypedHandler[T]) NeedsUpdate() bool {
	return th.facet.NeedsUpdate()
}

func (th *TypedHandler[T]) GetPageAny(first, pageSize int, sort sdk.SortSpec, filter string) (types.Page, error) {
	// This is where the magic happens - we use the type-safe facet methods
	// but return through the generic Page interface
	return th.GetPage(first, pageSize, sort, filter)
}

func (th *TypedHandler[T]) CrudAny(op crud.Operation, item interface{}) error {
	// Type assertion to convert interface{} back to T
	typedItem, ok := item.(T)
	if !ok {
		return fmt.Errorf("invalid type for %s operation: expected %T, got %T", th.dataFacet, *new(T), item)
	}
	return th.Crud(op, typedItem)
}

func (th *TypedHandler[T]) GetStoreName() string {
	return th.storeName
}

func (th *TypedHandler[T]) GetFacet() types.DataFacet {
	return th.dataFacet
}

// Type-safe methods that work with concrete types
func (th *TypedHandler[T]) GetPage(first, pageSize int, sort sdk.SortSpec, filter string) (types.Page, error) {
	// Create filter function from string
	var filterFunc func(*T) bool
	if filter != "" && th.filterFunc != nil {
		filterFunc = th.filterFunc.CreateFilter(filter)
	}

	// Create sort function
	var sortFunc func([]T, sdk.SortSpec) error
	if th.sortFunc != nil {
		sortFunc = th.sortFunc.CreateSort()
	}

	// Use the facet's GetPage method to get the paginated data
	pageResult, err := th.facet.GetPage(first, pageSize, filterFunc, sort, sortFunc)
	if err != nil {
		return nil, err
	}

	// Use the page factory to create the concrete page type
	if th.pageFactory != nil {
		return th.pageFactory.CreatePage(th.dataFacet, pageResult, th.facet)
	}

	return nil, fmt.Errorf("no page factory configured for %s", th.dataFacet)
}

func (th *TypedHandler[T]) Crud(op crud.Operation, item T) error {
	if th.crudHandler != nil {
		return th.crudHandler.Crud(op, item)
	}
	return fmt.Errorf("CRUD operation %s not supported for %s", op, th.dataFacet)
}

// GenericCollection provides a unified registry and interface for all collection operations
type GenericCollection struct {
	name     string
	handlers map[types.DataFacet]HandlerInterface
	mutex    sync.RWMutex
}

// NewGenericCollection creates a new generic collection with the given name
func NewGenericCollection(name string) *GenericCollection {
	return &GenericCollection{
		name:     name,
		handlers: make(map[types.DataFacet]HandlerInterface),
	}
}

// RegisterHandler registers a handler for a specific list dataFacet
func (gc *GenericCollection) RegisterHandler(dataFacet types.DataFacet, handler HandlerInterface) {
	gc.mutex.Lock()
	defer gc.mutex.Unlock()
	gc.handlers[dataFacet] = handler
}

// Implement the Collection interface
func (gc *GenericCollection) GetPage(dataFacet types.DataFacet, first, pageSize int, sort sdk.SortSpec, filter string) (types.Page, error) {
	gc.mutex.RLock()
	handler, exists := gc.handlers[dataFacet]
	gc.mutex.RUnlock()

	if !exists {
		return nil, fmt.Errorf("unsupported list dataFacet: %s", dataFacet)
	}

	return handler.GetPageAny(first, pageSize, sort, filter)
}

func (gc *GenericCollection) LoadData(dataFacet types.DataFacet) {
	gc.mutex.RLock()
	handler, exists := gc.handlers[dataFacet]
	gc.mutex.RUnlock()

	if exists {
		_ = handler.LoadData()
	}
}

func (gc *GenericCollection) Reset(dataFacet types.DataFacet) {
	gc.mutex.RLock()
	handler, exists := gc.handlers[dataFacet]
	gc.mutex.RUnlock()

	if exists {
		handler.Reset()
	}
}

func (gc *GenericCollection) NeedsUpdate(dataFacet types.DataFacet) bool {
	gc.mutex.RLock()
	handler, exists := gc.handlers[dataFacet]
	gc.mutex.RUnlock()

	if !exists {
		return false
	}

	return handler.NeedsUpdate()
}

func (gc *GenericCollection) Crud(dataFacet types.DataFacet, op crud.Operation, item interface{}) error {
	gc.mutex.RLock()
	handler, exists := gc.handlers[dataFacet]
	gc.mutex.RUnlock()

	if !exists {
		return fmt.Errorf("unsupported list dataFacet: %s", dataFacet)
	}

	return handler.CrudAny(op, item)
}

func (gc *GenericCollection) GetSupportedFacets() []types.DataFacet {
	gc.mutex.RLock()
	defer gc.mutex.RUnlock()

	dataFacets := make([]types.DataFacet, 0, len(gc.handlers))
	for dataFacet := range gc.handlers {
		dataFacets = append(dataFacets, dataFacet)
	}
	return dataFacets
}

func (gc *GenericCollection) GetStoreForFacet(dataFacet types.DataFacet) string {
	gc.mutex.RLock()
	handler, exists := gc.handlers[dataFacet]
	gc.mutex.RUnlock()

	if !exists {
		return ""
	}

	return handler.GetStoreName()
}

func (gc *GenericCollection) GetCollectionName() string {
	return gc.name
}
