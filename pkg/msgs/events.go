package msgs

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

var (
	wailsContext context.Context
	contextMutex sync.RWMutex
)

// Internal event bus for testing
var (
	testMode      bool
	testModeLock  sync.RWMutex
	listeners     map[EventType][]func(optionalData ...interface{})
	listenersLock sync.RWMutex
)

func init() {
	listeners = make(map[EventType][]func(optionalData ...interface{}))
}

func InitializeContext(ctx context.Context) {
	contextMutex.Lock()
	defer contextMutex.Unlock()
	wailsContext = ctx
	log.Println("Messaging context initialized")
}

func emitMessage(messageType EventType, msgText string, payload ...interface{}) {
	contextMutex.RLock()
	ctx := wailsContext
	contextMutex.RUnlock()

	if ctx != nil {
		// Create args slice with message first, then payload
		args := []interface{}{msgText}
		args = append(args, payload...)
		runtime.EventsEmit(ctx, string(messageType), args...)
	}

	if IsTestMode() {
		dispatchToListeners(messageType, msgText, payload...)
	}
}

// Sugar

func EmitLoaded(msgText string, payload ...interface{}) {
	emitMessage(EventDataLoaded, msgText, payload...)
}

func EmitStatus(msgText string, payload ...interface{}) {
	emitMessage(EventStatus, msgText, payload...)
}

func EmitManager(msgText string, payload ...interface{}) {
	emitMessage(EventManager, msgText, payload...)
}

func EmitError(msgText string, err error, payload ...interface{}) {
	msg := fmt.Sprintf("%s: %v", msgText, err)
	emitMessage(EventError, msg, payload...)
}

// SetTestMode enables or disables test mode
// In test mode, events are dispatched to registered listeners directly,
// bypassing the Wails runtime.
func SetTestMode(enabled bool) {
	testModeLock.Lock()
	defer testModeLock.Unlock()
	testMode = enabled
}

// IsTestMode returns whether test mode is enabled
func IsTestMode() bool {
	testModeLock.RLock()
	defer testModeLock.RUnlock()
	return testMode
}

// Internal function to register a listener for an event
func registerListener(eventType EventType, callback func(optionalData ...interface{})) func() {
	listenersLock.Lock()
	defer listenersLock.Unlock()

	// Create slice if it doesn't exist
	if _, exists := listeners[eventType]; !exists {
		listeners[eventType] = make([]func(optionalData ...interface{}), 0)
	}

	// Add the callback
	listeners[eventType] = append(listeners[eventType], callback)

	// Return an unsubscribe function
	return func() {
		listenersLock.Lock()
		defer listenersLock.Unlock()

		if listenersList, exists := listeners[eventType]; exists {
			// Find and remove the callback
			for i, registeredCallback := range listenersList {
				if &registeredCallback == &callback {
					// Remove by swapping with the last element and truncating
					lastIndex := len(listenersList) - 1
					if i != lastIndex {
						listenersList[i] = listenersList[lastIndex]
					}
					listeners[eventType] = listenersList[:lastIndex]
					break
				}
			}
		}
	}
}

// dispatchToListeners sends events to registered listeners in test mode
func dispatchToListeners(eventType EventType, msgText string, payload ...interface{}) {
	listenersLock.RLock()
	defer listenersLock.RUnlock()

	// Get listeners for this event type
	eventListeners, exists := listeners[eventType]
	if !exists || len(eventListeners) == 0 {
		return
	}

	// Create args slice with message first, then payload
	args := []interface{}{msgText}
	args = append(args, payload...)

	// Call each listener
	for _, listener := range eventListeners {
		// Call in a goroutine to prevent blocking
		go func(callback func(optionalData ...interface{})) {
			defer func() {
				if r := recover(); r != nil {
					log.Printf("Recovered from panic in event listener: %v", r)
				}
			}()
			callback(args...)
		}(listener)
	}
}

// On registers a listener for the specified event
// Returns a function to unregister the listener
func On(eventType EventType, callback func(optionalData ...interface{})) func() {
	// If in test mode, use internal event bus
	if IsTestMode() {
		return registerListener(eventType, callback)
	}

	// Otherwise, use Wails if context is available
	contextMutex.RLock()
	ctx := wailsContext
	contextMutex.RUnlock()

	if ctx != nil {
		return runtime.EventsOn(ctx, string(eventType), callback)
	}

	// Return a no-op unsubscribe function if context is not available
	return func() {}
}

func WaitForEvent(eventType EventType) <-chan bool {
	ch := make(chan bool, 1)

	// Register listener that will signal the channel
	unsub := On(eventType, func(optionalData ...interface{}) {
		// Only signal once and protect against double close
		select {
		case ch <- true:
			close(ch)
		default:
			// Channel already closed or full
		}
	})

	// Create a goroutine to clean up if the channel is never triggered
	go func() {
		<-ch
		unsub()
	}()

	return ch
}

func WaitForLoadedEvent(listKind string) <-chan bool {
	ch := make(chan bool, 1)

	// Register listener that will signal the channel if the list kind matches
	unsub := On(EventDataLoaded, func(optionalData ...interface{}) {
		if len(optionalData) > 1 {
			// Try to extract list kind from payload
			if payload, ok := optionalData[1].(map[string]interface{}); ok {
				if kind, exists := payload["listKind"]; exists && kind == listKind {
					// Only signal once and protect against double close
					select {
					case ch <- true:
						close(ch)
					default:
						// Channel already closed or full
					}
				}
			}
		}
	})

	// Create a goroutine to clean up if the channel is never triggered
	go func() {
		<-ch
		unsub()
	}()

	return ch
}
