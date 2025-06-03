package msgs

import (
	"context"
	"log"
	"sync"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

var (
	wailsContext context.Context
	contextMutex sync.RWMutex
)

func InitializeContext(ctx context.Context) {
	contextMutex.Lock()
	defer contextMutex.Unlock()
	wailsContext = ctx
	log.Println("Messaging context initialized")
}

// Update the EventEmitter type to not require context
type EventEmitter func(messageType EventType, msgText string, optionalPayload ...interface{})

// Update defaultEmitter to use stored context
var defaultEmitter EventEmitter = func(messageType EventType, msgText string, optionalPayload ...interface{}) {
	contextMutex.RLock()
	ctx := wailsContext
	contextMutex.RUnlock()

	if ctx != nil {
		runtime.EventsEmit(ctx, string(messageType), msgText, optionalPayload)
	}
}

func logEmitter(next EventEmitter) EventEmitter {
	return func(messageType EventType, msgText string, optionalPayload ...interface{}) {
		log.Printf("EVENT: %s - %s - %s", messageType, msgText, optionalPayload)
		next(messageType, msgText, optionalPayload...)
	}
}

var currentEmitter EventEmitter = defaultEmitter

func SetEmitter(emitter EventEmitter) {
	currentEmitter = emitter
}

func EnableLogging() {
	currentEmitter = logEmitter(defaultEmitter)
}

func DisableLogging() {
	currentEmitter = defaultEmitter
}

// EmitPayload emits a message using the stored context
func EmitPayload(messageType EventType, payload ...interface{}) {
	if currentEmitter != nil {
		currentEmitter(messageType, "payload", payload)
	}
}

// EmitMessage emits a message using the stored context
func EmitMessage(messageType EventType, msgText string, payload ...interface{}) {
	if currentEmitter != nil {
		currentEmitter(messageType, msgText, payload)
	}
}
