package source

import (
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/output"
)

// Source defines the contract for data sources
// Sources handle the low-level data fetching and streaming from external systems
type Source[T any] interface {
	Fetch(ctx *output.RenderCtx, processor func(item *T) bool) error
	GetSourceType() string
}

// SDKSource implements Source for TrueBlocks SDK-based data fetching
type SDKSource[T any] struct {
	queryFunc   func(*output.RenderCtx) error // SDK query function
	processFunc func(interface{}) *T          // Convert raw item to typed item
	sourceType  string                        // Type identifier for this source
}

// NewSDKSource creates a new SDK-based source
func NewSDKSource[T any](
	queryFunc func(*output.RenderCtx) error,
	processFunc func(interface{}) *T,
) *SDKSource[T] {
	return &SDKSource[T]{
		queryFunc:   queryFunc,
		processFunc: processFunc,
		sourceType:  "sdk",
	}
}

// Fetch implements Source.Fetch for SDK sources
func (s *SDKSource[T]) Fetch(ctx *output.RenderCtx, processor func(item *T) bool) error {
	// Set up done channel to coordinate fetching completion
	done := make(chan struct{})
	errChan := make(chan error, 1)

	// Start the SDK query in a goroutine
	go func() {
		defer func() {
			if ctx.ModelChan != nil {
				close(ctx.ModelChan)
			}
			if ctx.ErrorChan != nil {
				close(ctx.ErrorChan)
			}
			close(done)
		}()

		if err := s.queryFunc(ctx); err != nil {
			errChan <- err
		}
	}()

	// Process items as they come through the channels
	modelChanClosed := false
	errorChanClosed := false

	for !modelChanClosed || !errorChanClosed {
		select {
		case itemIntf, ok := <-ctx.ModelChan:
			if !ok {
				modelChanClosed = true
				if errorChanClosed {
					return nil
				}
				continue
			}

			// Convert the raw item to the expected type
			itemPtr := s.processFunc(itemIntf)
			if itemPtr == nil {
				continue // Skip items that couldn't be processed
			}

			// Process the item through the provided processor
			if !processor(itemPtr) {
				return nil
			}

		case streamErr, ok := <-ctx.ErrorChan:
			if !ok {
				errorChanClosed = true
				if modelChanClosed {
					return nil
				}
				continue
			}
			return streamErr

		case err := <-errChan:
			return err

		case <-done:
			// Streaming completed
			return nil

		case <-ctx.Ctx.Done():
			// Context was cancelled - this is the key addition!
			return ctx.Ctx.Err()
		}
	}

	return nil
}

// GetSourceType implements Source.GetSourceType
func (s *SDKSource[T]) GetSourceType() string {
	return s.sourceType
}
