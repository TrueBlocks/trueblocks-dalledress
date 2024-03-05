package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/file"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/logger"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/utils"
)

func main_images() {
	app := NewApp()

	ctx := context.Background()
	app.startup(ctx)
	app.domReady(ctx)
	// sourceFile := "addresses.txt"
	sourceFile := "blocks.txt"
	lines := file.AsciiFileToLines(sourceFile)
	if len(lines) > 0 {
		wg := sync.WaitGroup{}
		fn := "last.txt"
		v := strings.Trim(file.AsciiFileToString(fn), "\n\r")
		last := utils.MustParseInt(v)
		logger.Info("Starting at address ", last, " of ", len(lines), file.FileExists(fn), v)
		for i := 0; i < len(lines); i++ {
			if lines[i][0] == '#' || len(lines[i]) < 42 {
				continue
			}
			if i > int(last) {
				if address, ok := app.validateInput(lines[i]); !ok {
					fmt.Println("Invalid address", lines[i])
					return
				} else {
					wg.Add(1)
					go doOne(i, &wg, app, address.Hex())
					file.StringToAsciiFile(fn, fmt.Sprintf("%d\n", i))
					if (i+1)%5 == 0 {
						wg.Wait()
						logger.Info("Sleeping for 60 seconds")
						time.Sleep(time.Second * 60)
					}
				}
			}
		}
		wg.Wait()
		return
	}

	if len(os.Args) == 2 {
		for i := 0; i < 10; i++ {
			wg := sync.WaitGroup{}
			for j := 0; j < 5; j++ {
				logger.Info("Round", i, "run", j)
				wg.Add(1)
				go doOne(i, &wg, app, fmt.Sprintf("0x%040x", 10010010+(i*10)+j)) // os.Args[1])
				// SeedBump++
			}
			wg.Wait()
			logger.Info("Sleeping for 60 seconds")
			time.Sleep(time.Second * 60)
		}
	} else {
		wg := sync.WaitGroup{}
		for i, arg := range os.Args[1:] {
			wg.Add(1)
			go doOne(i, &wg, app, arg)
		}
		wg.Wait()
	}
}
