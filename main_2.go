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
	lines := file.AsciiFileToLines("addresses.txt")
	if len(lines) > 0 {
		wg := sync.WaitGroup{}
		for i := 0; i < len(lines); i++ {
			// fmt.Println(lines[i])
			wg.Add(1)
			go doOne(&wg, app, lines[i])
			if (i+1)%5 == 0 {
				wg.Wait()
				logger.Info("Sleeping for 60 seconds")
				time.Sleep(time.Second * 60)
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
				go doOne(&wg, app, fmt.Sprintf("0x%040x", 10010010+(i*10)+j)) // os.Args[1])
				// SeedBump++
			}
			wg.Wait()
			logger.Info("Sleeping for 60 seconds")
			time.Sleep(time.Second * 60)
		}
	} else {
		wg := sync.WaitGroup{}
		for _, arg := range os.Args[1:] {
			wg.Add(1)
			go doOne(&wg, app, arg)
		}
		wg.Wait()
	}
}
