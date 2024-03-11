package main

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/file"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/logger"
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
		logger.Info("Starting at address ", app.Series.Last, " of ", len(lines))
		app.nMade = 0
		for i := 0; i < len(lines); i++ {
			if lines[i][0] == '#' || len(lines[i]) < 42 {
				continue
			}
			if i > int(app.Series.Last) {
				if address, ok := app.validateInput(lines[i]); !ok {
					fmt.Println("Invalid address", lines[i])
					return
				} else {
					wg.Add(1)
					go doOne(i, &wg, app, address.Hex())
					app.Series.Last = i
					app.Series.Save()
					if (i+1)%5 == 0 {
						wg.Wait()
						if app.nMade > 4 {
							logger.Info("Sleeping for 60 seconds")
							time.Sleep(time.Second * 60)
							app.nMade = 0
						}
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
