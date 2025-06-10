package store

// // TestReloadCancellation tests that Reload properly cancels ongoing operations
// func TestReloadCancellation(t *testing.T) {
// 	abisAddr := base.ZeroAddr.Hex()
// 	renderCtx := RegisterContext(abisAddr)

// 	// Verify the context was registered
// 	if ctxCount(abisAddr) != 1 {
// 		t.Errorf("Expected 1 registered context, got %d", ctxCount(abisAddr))
// 	}

// 	// Verify the context is not nil
// 	if renderCtx == nil {
// 		t.Error("RegisterContext should return non-nil context")
// 	}

// 	// Simulate a reload operation by cancelling the context
// 	cancelled, found := UnregisterContext(abisAddr)
// 	if !found {
// 		t.Error("Cancel should find the registered context")
// 	}
// 	if cancelled != 1 {
// 		t.Errorf("Expected 1 cancelled context, got %d", cancelled)
// 	}

// 	// Verify the context was cancelled and removed
// 	if ctxCount(abisAddr) != 0 {
// 		t.Errorf("Expected 0 registered contexts after reload, got %d", ctxCount(abisAddr))
// 	}

// 	// Note: We can't easily test if the context was actually cancelled since
// 	// the Cancel method removes it from the map, but the fact that it was
// 	// removed indicates it was processed correctly
// }

// // TestContextRegistration tests that contexts are properly registered and cleaned up
// func TestContextRegistration(t *testing.T) {
// 	addr1 := "0x1234567890123456789012345678901234567890"
// 	addr2 := "0x2234567890123456789012345678901234567890"

// 	ctx1 := RegisterContext(addr1)
// 	ctx2 := RegisterContext(addr2)

// 	cnt := ctxCount(addr1) + ctxCount(addr1)
// 	if cnt != 2 {
// 		t.Errorf("Expected 2 registered contexts, got %d", cnt)
// 	}

// 	if ctx1 == nil || ctx2 == nil {
// 		t.Error("RegisterContext should return non-nil contexts")
// 	}

// 	// Test Cancel for specific address
// 	cancelled, found := UnregisterContext(addr1)
// 	if !found {
// 		t.Error("Cancel should find the registered context")
// 	}
// 	if cancelled != 1 {
// 		t.Errorf("Expected 1 cancelled context, got %d", cancelled)
// 	}
// 	cnt = ctxCount(addr1) + ctxCount(addr2)
// 	if cnt != 1 {
// 		t.Errorf("Expected 1 remaining context after cancel, got %d", cnt)
// 	}

// 	// Test Cancel for non-existent address
// 	nonExistentAddr := "0x9999999999999999999999999999999999999999"
// 	cancelled, found = UnregisterContext(nonExistentAddr)
// 	if found {
// 		t.Error("Cancel should not find non-existent context")
// 	}
// 	if cancelled != 0 {
// 		t.Errorf("Expected 0 cancelled contexts for non-existent address, got %d", cancelled)
// 	}
// }
