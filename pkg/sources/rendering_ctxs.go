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
	renderCtxs      map[string][]*output.RenderCtx
	renderCtxsMutex sync.Mutex
}

// GetContextManager returns the singleton context manager
func GetContextManager() *ContextManager {
	initOnce.Do(func() {
		globalContextManager = &ContextManager{
			renderCtxs: make(map[string][]*output.RenderCtx),
		}
	})
	return globalContextManager
}

// RegisterCtx registers a new RenderCtx for a given key
func RegisterCtx(key string) *output.RenderCtx {
	cm := GetContextManager()
	cm.renderCtxsMutex.Lock()
	defer cm.renderCtxsMutex.Unlock()

	rCtx := output.NewStreamingContext()
	cm.renderCtxs[key] = append(cm.renderCtxs[key], rCtx)
	return rCtx
}

// CtxCount returns the number of contexts for a key
func CtxCount(key string) int {
	cm := GetContextManager()
	cm.renderCtxsMutex.Lock()
	defer cm.renderCtxsMutex.Unlock()

	return len(cm.renderCtxs[key])
}

// Cancel closes and removes all contexts for a key
func Cancel(key string) (int, bool) {
	cm := GetContextManager()
	cm.renderCtxsMutex.Lock()
	defer cm.renderCtxsMutex.Unlock()

	if len(cm.renderCtxs) == 0 {
		return 0, false
	}
	if cm.renderCtxs[key] == nil {
		return 0, false
	}
	n := len(cm.renderCtxs[key])
	for i := 0; i < len(cm.renderCtxs[key]); i++ {
		cm.renderCtxs[key][i].Cancel()
	}
	cm.renderCtxs[key] = nil
	return n, true
}

func CancelAll() {
	cm := GetContextManager()
	cm.renderCtxsMutex.Lock()
	defer cm.renderCtxsMutex.Unlock()

	for key := range cm.renderCtxs {
		if cm.renderCtxs[key] != nil {
			for i := 0; i < len(cm.renderCtxs[key]); i++ {
				cm.renderCtxs[key][i].Cancel()
			}
			cm.renderCtxs[key] = nil
		}
	}
}
