package app

import (
	"path/filepath"

	"github.com/TrueBlocks/trueblocks-dalledress/pkg/utils"
)

func (a *App) GetFilelist() []string {
	folder := filepath.Join("./output/series")
	if list, err := utils.Listing(folder); err != nil {
		return []string{err.Error()}
	} else {
		return list
	}
}
