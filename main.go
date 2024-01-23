/*
package main

import (
	"context"
)

func main() {
	app := NewApp()
	app.startup(context.Background())
	app.GetImage("trueblocks.eth")
}

*/

package main

import (
	"context"
	"embed"
	"fmt"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/file"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// Create an instance of the app structure
	app := NewApp()
	app.startup(context.Background())

	x := true
	if x {
		// Create application with options
		err := wails.Run(&options.App{
			Title:  "dalledresses",
			Width:  1024,
			Height: 768,
			AssetServer: &assetserver.Options{
				Assets: assets,
			},
			BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
			OnStartup:        app.startup,
			Bind: []interface{}{
				app,
				// &SimpleBounds{},
			},
		})

		if err != nil {
			println("Error:", err.Error())
		}
	} else {
		fn := "/Users/jrush/Desktop/Animals.1/addresses2.csv"
		addrs := file.AsciiFileToLines(fn)
		// lens := []int{len(app.adverbs), len(app.adjectives), len(app.nouns), len(app.styles), len(app.colors), len(app.colors), len(app.colors), len(app.styles)}
		fmt.Println("[")
		for index, addr := range addrs {
			if index > 0 {
				fmt.Println(",")
			}
			// fmt.Println(app.GetData(addr))
			// fmt.Println(strings.Repeat("-", 80))
			fmt.Println(app.GetJson(addr))
		}
		fmt.Println("]")
	}
}
