package app

import (
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/names"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/types"
)

// TODO: Eventually this will get put back into Core.

type NameEx struct {
	types.Name `json:",inline"`
	Type       names.Parts `json:"type"`
}

var NameDbParts = []struct {
	Value  names.Parts
	TSName string
}{
	{names.Regular, "REGULAR"},
	{names.Custom, "CUSTOM"},
	{names.Prefund, "PREFUND"},
	{names.Baddress, "BADDRESS"},
}
