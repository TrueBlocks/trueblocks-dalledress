package types

import (
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/output"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/msgs"
)

type App interface {
	EmitEvent(eventType msgs.EventType, payload interface{})
	// TODO: This is incorrect. We need to register the rendering context
	// TODO: in the app so we can cancel it later. I'm not sure this does that.
	RegisterCtx(key string) *output.RenderCtx
	Cancel(key string) (int, bool)
	LogBackend(msg string)
}
