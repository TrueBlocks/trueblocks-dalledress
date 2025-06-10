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

func SetTestMode(enabled bool) {
	testModeLock.Lock()
	defer testModeLock.Unlock()
	testMode = enabled
}

func IsTestMode() bool {
	testModeLock.RLock()
	defer testModeLock.RUnlock()
	return testMode
}

func registerListener(eventType EventType, callback func(optionalData ...interface{})) func() {
	listenersLock.Lock()
	defer listenersLock.Unlock()

	if _, exists := listeners[eventType]; !exists {
		listeners[eventType] = make([]func(optionalData ...interface{}), 0)
	}

	listeners[eventType] = append(listeners[eventType], callback)
	return func() {
		listenersLock.Lock()
		defer listenersLock.Unlock()
		if listenersList, exists := listeners[eventType]; exists {
			for i, registeredCallback := range listenersList {
				if &registeredCallback == &callback {
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

func dispatchToListeners(eventType EventType, msgText string, payload ...interface{}) {
	listenersLock.RLock()
	defer listenersLock.RUnlock()

	eventListeners, exists := listeners[eventType]
	if !exists || len(eventListeners) == 0 {
		return
	}

	args := []interface{}{msgText}
	args = append(args, payload...)

	for _, listener := range eventListeners {
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

func On(eventType EventType, callback func(optionalData ...interface{})) func() {
	if IsTestMode() {
		return registerListener(eventType, callback)
	}

	contextMutex.RLock()
	ctx := wailsContext
	contextMutex.RUnlock()

	if ctx != nil {
		return runtime.EventsOn(ctx, string(eventType), callback)
	}

	return func() {}
}

func WaitForEvent(eventType EventType) <-chan bool {
	ch := make(chan bool, 1)

	unsub := On(eventType, func(optionalData ...interface{}) {
		select {
		case ch <- true:
			close(ch)
		default:
			// Channel already closed or full
		}
	})

	go func() {
		<-ch
		unsub()
	}()

	return ch
}

func WaitForLoadedEvent(listKind string) <-chan bool {
	ch := make(chan bool, 1)
	unsub := On(EventDataLoaded, func(optionalData ...interface{}) {
		if len(optionalData) > 1 {
			if payload, ok := optionalData[1].(map[string]interface{}); ok {
				if kind, exists := payload["listKind"]; exists && kind == listKind {
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

	go func() {
		<-ch
		unsub()
	}()

	return ch
}
