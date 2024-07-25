package app

import (
	"path/filepath"

	"github.com/TrueBlocks/trueblocks-dalledress/pkg/utils"
)

func (a *App) GetFilelist(baseFolder string) []string {
	folder := filepath.Join(baseFolder)
	if list, err := utils.Listing(folder); err != nil {
		return []string{err.Error()}
	} else {
		return list
	}
}
