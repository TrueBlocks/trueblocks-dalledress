package app

import (
	"testing"

	"github.com/TrueBlocks/trueblocks-dalledress/pkg/types"
	daltypes "github.com/TrueBlocks/trueblocks-dalledress/pkg/types/dalledress"
	sdk "github.com/TrueBlocks/trueblocks-sdk/v5"
	"github.com/stretchr/testify/require"
)

// Test that GetDalleDressPage returns a page with non-nil CurrentDress for generator facet
func TestGetDalleDressPageCurrentDress(t *testing.T) {
	a := &App{}
	payload := &types.Payload{DataFacet: daltypes.DalleDressGenerator}
	var sortSpec sdk.SortSpec
	page, err := a.GetDalleDressPage(payload, 0, 10, sortSpec, "")
	require.NoError(t, err)
	require.NotNil(t, page)
	require.Equal(t, daltypes.DalleDressGenerator, page.Facet)
	require.NotNil(t, page.CurrentDress)
}
