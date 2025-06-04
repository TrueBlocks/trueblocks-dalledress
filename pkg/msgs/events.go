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

func EmitPayload(messageType EventType, payload interface{}) {
	contextMutex.RLock()
	ctx := wailsContext
	contextMutex.RUnlock()
	if ctx != nil {
		runtime.EventsEmit(ctx, string(messageType), payload)
	}
}

func EmitMessage(messageType EventType, msgText string, payload ...interface{}) {
	contextMutex.RLock()
	ctx := wailsContext
	contextMutex.RUnlock()

	if ctx != nil {
		runtime.EventsEmit(ctx, string(messageType), payload...)
	}
}
