package app

import (
	"testing"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/base"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/types"
	"github.com/stretchr/testify/assert"
)

func setupTestApp() *App {
	app := &App{}
	// Setup minimal test
	app.names.Map = make(map[base.Address]types.Name)
	app.names.List = []*types.Name{}

	return app
}

func TestModifyName_Add(t *testing.T) {
	// This test is more of an integration test and requires the SDK
	// We'll mock/skip the actual SDK call but test the local state updates

	t.Skip("Skipping integration test that requires SDK")

	app := setupTestApp()
	addr := base.HexToAddress("0x123456789abcdef0123456789abcdef012345678")

	// Add a new name
	err := app.ModifyName(&ModifyData{
		Operation: "add",
		Address:   addr,
		Value:     "TestName",
	})

	// In real test, check err is nil
	assert.Nil(t, err)

	// Check if name was added to local state
	name, exists := app.names.Map[addr]
	assert.True(t, exists)
	assert.Equal(t, "TestName", name.Name)
	assert.True(t, name.IsCustom)
}

func TestModifyName_Edit(t *testing.T) {
	t.Skip("Skipping integration test that requires SDK")

	app := setupTestApp()
	addr := base.HexToAddress("0x123456789abcdef0123456789abcdef012345678")

	// First add a name
	app.names.Map[addr] = types.Name{
		Address:  addr,
		Name:     "OldName",
		IsCustom: true,
		Tags:     "99-User-Defined",
	}
	name := app.names.Map[addr]
	app.names.List = append(app.names.List, &name)

	// Then edit it
	err := app.ModifyName(&ModifyData{
		Operation: "update",
		Address:   addr,
		Value:     "NewName",
	})

	// In real test, check err is nil
	assert.Nil(t, err)

	// Check if name was updated in local state
	name, exists := app.names.Map[addr]
	assert.True(t, exists)
	assert.Equal(t, "NewName", name.Name)
	assert.True(t, name.IsCustom)
}

func TestModifyName_Delete(t *testing.T) {
	t.Skip("Skipping integration test that requires SDK")

	app := setupTestApp()
	addr := base.HexToAddress("0x123456789abcdef0123456789abcdef012345678")

	// First add a name
	app.names.Map[addr] = types.Name{
		Address:  addr,
		Name:     "ToDelete",
		IsCustom: true,
		Tags:     "99-User-Defined",
	}
	name := app.names.Map[addr]
	app.names.List = append(app.names.List, &name)

	// Then delete it
	err := app.ModifyName(&ModifyData{
		Operation: "delete",
		Address:   addr,
		Value:     "",
	})

	// In real test, check err is nil
	assert.Nil(t, err)

	// Check if name was marked as deleted in local state
	name, exists := app.names.Map[addr]
	assert.True(t, exists)
	assert.True(t, name.Deleted)
}

func TestModifyName_Undelete(t *testing.T) {
	t.Skip("Skipping integration test that requires SDK")

	app := setupTestApp()
	addr := base.HexToAddress("0x123456789abcdef0123456789abcdef012345678")

	// First add a deleted name
	app.names.Map[addr] = types.Name{
		Address:  addr,
		Name:     "ToUndelete",
		IsCustom: true,
		Tags:     "99-User-Defined",
		Deleted:  true,
	}
	name := app.names.Map[addr]
	app.names.List = append(app.names.List, &name)

	// Then undelete it
	err := app.ModifyName(&ModifyData{
		Operation: "undelete",
		Address:   addr,
		Value:     "",
	})

	// In real test, check err is nil
	assert.Nil(t, err)

	// Check if name was marked as not deleted in local state
	name, exists := app.names.Map[addr]
	assert.True(t, exists)
	assert.False(t, name.Deleted)
}

func TestModifyName_Remove(t *testing.T) {
	t.Skip("Skipping integration test that requires SDK")

	app := setupTestApp()
	addr := base.HexToAddress("0x123456789abcdef0123456789abcdef012345678")

	// First add a name
	app.names.Map[addr] = types.Name{
		Address:  addr,
		Name:     "ToRemove",
		IsCustom: true,
		Tags:     "99-User-Defined",
	}
	name := app.names.Map[addr]
	app.names.List = append(app.names.List, &name)

	// Then remove it
	err := app.ModifyName(&ModifyData{
		Operation: "remove",
		Address:   addr,
		Value:     "",
	})

	// In real test, check err is nil
	assert.Nil(t, err)

	// Check if name was removed from local state
	_, exists := app.names.Map[addr]
	assert.False(t, exists)

	// Check if name was removed from list
	found := false
	for _, name := range app.names.List {
		if name.Address == addr {
			found = true
			break
		}
	}
	assert.False(t, found)
}

func TestModifyName_Locking(t *testing.T) {
	app := setupTestApp()
	addr := base.HexToAddress("0x123456789abcdef0123456789abcdef012345678")

	// Set the lock manually
	namesLock.CompareAndSwap(0, 1)

	// Try to modify a name while locked
	err := app.ModifyName(&ModifyData{
		Operation: "add",
		Address:   addr,
		Value:     "TestName",
	})

	// Should return nil because the function exits early when locked
	assert.Nil(t, err)

	// Name should not be added because operation was blocked by lock
	_, exists := app.names.Map[addr]
	assert.False(t, exists)

	// Reset the lock for other tests
	namesLock.CompareAndSwap(1, 0)
}
