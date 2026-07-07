package main

import (
	"embed"
	"log"
	"os"
	"path/filepath"

	appkit "github.com/TrueBlocks/trueblocks-art/packages/appkit/v2"
	"github.com/TrueBlocks/trueblocks-dalledress/v2/app"
	"github.com/TrueBlocks/trueblocks-dalledress/v2/internal/db"
	"github.com/wailsapp/wails/v2/pkg/options"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	dataDir := appkit.AppDirFor("dalledress")
	if err := os.MkdirAll(dataDir, appkit.DirPermissions); err != nil {
		log.Fatalf("data dir: %v", err)
	}

	database, err := db.Open(filepath.Join(dataDir, "dalledress.db"))
	if err != nil {
		log.Fatalf("open db: %v", err)
	}
	defer func() { _ = database.Close() }()

	application, err := app.NewApp(filepath.Join(dataDir, "prefs.json"), database)
	if err != nil {
		log.Fatalf("new app: %v", err)
	}

	err = appkit.Run(appkit.AppConfig{
		Title:             "Dalledress",
		Assets:            assets,
		BackgroundColour:  &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:         application.Startup,
		OnShutdown:        application.Shutdown,
		GetWindowGeometry: application.GetWindowGeometry,
		Bind: []any{
			application,
		},
	})
	if err != nil {
		log.Fatal(err)
	}
}
