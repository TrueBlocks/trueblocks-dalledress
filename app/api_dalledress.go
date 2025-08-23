// Copyright 2016, 2026 The Authors. All rights reserved.
// Use of this source code is governed by a license that can
// be found in the LICENSE file.
/*
 * Parts of this file were auto generated. Edit only those parts of
 * the code inside of 'EXISTING_CODE' tags.
 */

package app

import (
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/types"
)

// GetDalleDressConfig returns the view configuration for dalledress
func (a *App) GetDalleDressConfig(payload types.Payload) (*types.ViewConfig, error) {
	return &types.ViewConfig{
		FacetOrder: []string{"generator"},
		Facets: map[string]types.FacetConfig{
			"generator": {
				Name:          "Generator",
				Store:         "generator",
				IsForm:        true,
				DividerBefore: false,
				Fields:        []types.FieldConfig{},
				Actions:       []string{},
				HeaderActions: []string{},
				RendererTypes: "",
			},
		},
	}, nil
}
