package sorting

import sdk "github.com/TrueBlocks/trueblocks-sdk/v5"

type SortDef struct {
	Key       string `json:"key"`
	Direction string `json:"direction"`
}

// ConvertToSortSpec converts our SortDef to SDK's SortSpec format
func ConvertToSortSpec(sortDef *SortDef) sdk.SortSpec {
	if sortDef == nil || sortDef.Key == "" {
		return sdk.SortSpec{}
	}

	order := sdk.Asc
	if sortDef.Direction == "desc" {
		order = sdk.Dec
	}

	return sdk.SortSpec{
		Fields: []string{sortDef.Key},
		Order:  []sdk.SortOrder{order},
	}
}
