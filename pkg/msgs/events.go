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
type EventEmitter func(messageType EventType, msgText string)

// Update defaultEmitter to use stored context
var defaultEmitter EventEmitter = func(messageType EventType, msgText string) {
	contextMutex.RLock()
	ctx := wailsContext
	contextMutex.RUnlock()

	if ctx != nil {
		runtime.EventsEmit(ctx, string(messageType), msgText)
	}
}

func logEmitter(next EventEmitter) EventEmitter {
	return func(messageType EventType, msgText string) {
		log.Printf("EVENT: %s - %s", messageType, msgText)
		next(messageType, msgText)
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

// EmitMessage emits a message using the stored context
func EmitMessage(messageType EventType, msgText string) {
	if currentEmitter != nil {
		currentEmitter(messageType, msgText)
	}
}
