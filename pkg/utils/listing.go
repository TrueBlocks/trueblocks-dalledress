package utils

import (
	"context"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/colors"
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
			ret = append(ret, result.Path)
			logger.Info(colors.Green, result.Path, colors.Off)
			// ok, err := walker.visitFunc1(walker, result.Path, cnt == 0)
			// if err != nil {
			// 	return err
			// }
			// if ok {
			cnt++
			// } else {
			// 	return nil
			// }
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
