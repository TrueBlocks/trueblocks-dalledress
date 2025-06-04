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
