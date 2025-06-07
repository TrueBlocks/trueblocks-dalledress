package sources

import (
	"sync"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/output"
)

var (
	globalContextManager *ContextManager
	initOnce             sync.Once
)

type ContextManager struct {
	renderCtxs      map[string]*output.RenderCtx
	renderCtxsMutex sync.Mutex
}

// GetContextManager returns the singleton context manager
func GetContextManager() *ContextManager {
	initOnce.Do(func() {
		globalContextManager = &ContextManager{
			renderCtxs: make(map[string]*output.RenderCtx),
		}
	})
	return globalContextManager
}

// RegisterContext registers a new RenderCtx for a given key
func RegisterContext(key string) *output.RenderCtx {
	cm := GetContextManager()
	cm.renderCtxsMutex.Lock()
	defer cm.renderCtxsMutex.Unlock()

	if existingCtx := cm.renderCtxs[key]; existingCtx != nil {
		existingCtx.Cancel()
	}

	rCtx := output.NewStreamingContext()
	cm.renderCtxs[key] = rCtx
	return rCtx
}

// UnregisterContext removes and cancels the context for a given key
func UnregisterContext(key string) (int, bool) {
	cm := GetContextManager()
	cm.renderCtxsMutex.Lock()
	defer cm.renderCtxsMutex.Unlock()

	if len(cm.renderCtxs) == 0 {
		return 0, false
	}
	if cm.renderCtxs[key] == nil {
		return 0, false
	}

	cm.renderCtxs[key].Cancel()
	delete(cm.renderCtxs, key)
	return 1, true
}

// CtxCount returns the number of contexts for a key
func CtxCount(key string) int {
	cm := GetContextManager()
	cm.renderCtxsMutex.Lock()
	defer cm.renderCtxsMutex.Unlock()

	if cm.renderCtxs[key] != nil {
		return 1
	}
	return 0
}

func CancelFetch(contextKey string) {
	cm := GetContextManager()
	cm.renderCtxsMutex.Lock()
	defer cm.renderCtxsMutex.Unlock()

	for key, ctx := range cm.renderCtxs {
		if key == contextKey && ctx != nil {
			ctx.Cancel()
			delete(cm.renderCtxs, key)
		}
	}
}
