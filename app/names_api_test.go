package app

import (
	"testing"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/base"
	coreTypes "github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/types"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/preferences"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/types"
	"github.com/stretchr/testify/assert"
)

func setupTestApp() *App {
	app := &App{
		Preferences: &preferences.Preferences{
			Org:  preferences.OrgPreferences{},
			User: preferences.UserPreferences{},
			App:  preferences.AppPreferences{},
		},
	}
	// ADD_ABIS_CODE
	app.names = types.NewNamesCollection(app)
	app.abis = types.NewAbisCollection(app)
	_ = app.Reload()
	return app
}

func TestModifyName(t *testing.T) {
	app := setupTestApp()
	name := &coreTypes.Name{
		Address: base.HexToAddress("0xbe0ed4138121ecfc5c0e56b40517da27e6c5226b"),
		Name:    "TestName",
	}

	// tests some of the fields for equality
	equals := func(t *testing.T, name string, old *coreTypes.Name, new *coreTypes.Name) bool {
		assert.Equal(t, name, new.Name)
		assert.True(t, old.Address == new.Address)
		assert.True(t, old.Tags == new.Tags)
		assert.True(t, old.Source == new.Source)
		assert.True(t, old.Decimals == new.Decimals)
		assert.Equal(t, coreTypes.Custom, new.Parts)
		assert.True(t, base.ZeroWei.Cmp(&old.Prefund) == 0)
		assert.True(t, new.IsCustom)
		return !t.Failed()
	}

	t.Run("Add", func(t *testing.T) {
		err := app.ModifyName("create", name)
		assert.Nil(t, err)
		n, exists := app.names.Map[name.Address]
		assert.True(t, exists)
		assert.True(t, equals(t, "TestName", name, &n))
		assert.False(t, name.Deleted)
		assert.False(t, name.IsContract)
		assert.False(t, name.IsErc20)
		assert.False(t, name.IsErc721)
		assert.False(t, name.IsPrefund)
	})

	t.Run("Autoname", func(t *testing.T) {
		err := app.ModifyName("autoname", name)
		assert.Nil(t, err)
		n, exists := app.names.Map[name.Address]
		assert.True(t, exists)
		assert.False(t, n.Deleted)
		name = &n
	})

	t.Run("Edit", func(t *testing.T) {
		name.Name = "UpdatedName"
		err := app.ModifyName("update", name)
		assert.Nil(t, err)
		n, exists := app.names.Map[name.Address]
		assert.True(t, exists)
		assert.Equal(t, "UpdatedName", n.Name)
		assert.True(t, n.IsCustom)
	})

	t.Run("Delete", func(t *testing.T) {
		err := app.ModifyName("delete", name)
		assert.Nil(t, err)
		n, exists := app.names.Map[name.Address]
		assert.True(t, exists)
		assert.True(t, n.Deleted)
	})

	t.Run("Undelete", func(t *testing.T) {
		err := app.ModifyName("undelete", name)
		assert.Nil(t, err)
		n, exists := app.names.Map[name.Address]
		assert.True(t, exists)
		assert.False(t, n.Deleted)
	})

	t.Run("Remove not deleted will fail", func(t *testing.T) {
		err := app.ModifyName("remove", name)
		assert.NotNil(t, err)
		n, exists := app.names.Map[name.Address]
		assert.True(t, exists)
		assert.False(t, n.Deleted)
	})

	t.Run("Delete again", func(t *testing.T) {
		err := app.ModifyName("delete", name)
		assert.Nil(t, err)
		n, exists := app.names.Map[name.Address]
		assert.True(t, exists)
		assert.True(t, n.Deleted)
	})

	t.Run("Remove", func(t *testing.T) {
		err := app.ModifyName("remove", name)
		assert.Nil(t, err)
		_, exists := app.names.Map[name.Address]
		assert.False(t, exists)
		found := false
		for _, nn := range app.names.List {
			if nn.Address == name.Address {
				found = true
				break
			}
		}
		assert.False(t, found)
	})

	t.Run("Locking", func(t *testing.T) {
		namesLock.CompareAndSwap(0, 1)
		err := app.ModifyName("add", name)
		assert.Nil(t, err)
		_, exists := app.names.Map[name.Address]
		assert.False(t, exists)
		namesLock.CompareAndSwap(1, 0)
	})
}
