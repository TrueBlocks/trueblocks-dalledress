package utils

import (
	"context"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/file"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/logger"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/walk"
)

func Listing(folderPath string) ([]string, error) {
	ret := []string{}
	filenameChan := make(chan walk.CacheFileInfo)

	var nRoutines int = 1
	go walk.WalkFolder(context.Background(), folderPath, nil, filenameChan)

	cnt := 0
	for result := range filenameChan {
		switch result.Type {
		case walk.Regular:
			if file.FileExists(result.Path) {
				ret = append(ret, result.Path)
				cnt++
			}
		case walk.Cache_NotACache:
			nRoutines--
			if nRoutines == 0 {
				close(filenameChan)
				continue
			}
		default:
			logger.Fatal("should not happen ==> in WalkRegularFolder")
		}
	}

	return ret, nil
}
