package app

import (
	"path/filepath"

	"github.com/TrueBlocks/trueblocks-dalledress/pkg/utils"
)

func (a *App) GetSeries(baseFolder string) []utils.Series {
	folder := filepath.Join(baseFolder)
	if list, err := utils.Listing(folder); err != nil {
		return []utils.Series{utils.Series(err.Error())}
	} else {
		return list
	}
}
