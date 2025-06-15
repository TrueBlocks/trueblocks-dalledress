package store

import (
	"testing"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/base"
	"github.com/stretchr/testify/assert"
)

func TestReloadCancellation(t *testing.T) {
	abisAddr := base.ZeroAddr.Hex()
	renderCtx := RegisterContext(abisAddr)

	assert.Equal(t, 1, ctxCount(abisAddr), "Expected 1 registered context")
	assert.NotNil(t, renderCtx, "RegisterContext should return non-nil context")

	cancelled, found := UnregisterContext(abisAddr)
	assert.True(t, found, "UnregisterContext should find the registered context")
	assert.Equal(t, 1, cancelled, "Expected 1 cancelled context")
	assert.Equal(t, 0, ctxCount(abisAddr), "Expected 0 registered contexts after reload")

	// Note: We can't easily test if the context was actually cancelled since
	// the Cancel method removes it from the map, but the fact that it was
	// removed indicates it was processed correctly
}

func TestContextRegistration(t *testing.T) {
	addr1 := "0x1234567890123456789012345678901234567890"
	addr2 := "0x2234567890123456789012345678901234567890"

	ctx1 := RegisterContext(addr1)
	ctx2 := RegisterContext(addr2)

	assert.True(t, ctxCount(addr1)+ctxCount(addr1) == 2, "Expected 2 registered contexts")
	assert.True(t, ctx1 != nil && ctx2 != nil, "Contexts should not be nil after registration")

	cancelled, found := UnregisterContext(addr1)
	assert.True(t, found, "UnregisterContext should find the registered context")
	assert.True(t, cancelled == 1, "Expected 1 cancelled context")
	assert.True(t, ctxCount(addr1)+ctxCount(addr2) == 1, "Expected 1 remaining context after cancellation")

	nonExistentAddr := "0x9999999999999999999999999999999999999999"
	cancelled, found = UnregisterContext(nonExistentAddr)
	assert.False(t, found, "UnregisterContext should not find non-existent context")
	assert.True(t, cancelled == 0, "Expected 0 cancelled contexts for non-existent address")
}
