package main

import (
	"context"
	"embed"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/logger"
	"github.com/TrueBlocks/trueblocks-dalledress/app"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/messages"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/types"
	"github.com/TrueBlocks/trueblocks-dalledress/servers"
	"github.com/wailsapp/wails/v2"
	wLogger "github.com/wailsapp/wails/v2/pkg/logger"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	if os.Getenv("TB_CMD_LINE") == "true" {
		logger.Info("Running in console mode")
		a := app.NewApp()
		ctx := context.Background()
		a.Startup(ctx)
		a.DomReady(ctx)
		a.HandleLines()
	} else {
		a := app.NewApp()
		opts := options.App{
			Title:            a.GetSession().Title,
			Width:            a.GetSession().Width,
			Height:           a.GetSession().Height,
			OnStartup:        a.Startup,
			OnDomReady:       a.DomReady,
			OnShutdown:       a.Shutdown,
			BackgroundColour: nil,
			LogLevel:         wLogger.ERROR,
			Menu:             a.GetMenus(),
			// Find: NewViews
			Bind: []interface{}{
				a,
				&messages.DocumentMsg{},
				&messages.ServerMsg{},
				&messages.ProgressMsg{},
				&messages.ErrorMsg{},
				&types.NameEx{},
				&types.TransactionEx{},
				&servers.Server{},
				&types.MonitorEx{},
			},
			EnumBind: []interface{}{
				types.NameParts,
				servers.Types,
				servers.States,
				messages.Messages,
			},
			StartHidden: true,
			AssetServer: &assetserver.Options{
				Assets: assets,
			},
		}

		http.HandleFunc("/files/", func(w http.ResponseWriter, r *http.Request) {
			address := strings.TrimPrefix(r.URL.Path, "/files/")
			parts := strings.Split(address, "&")
			if len(parts) > 1 {
				address = parts[0]
			}
			if address == "" {
				http.Error(w, "Address not provided", http.StatusBadRequest)
				return
			}
			cwd, err := os.Getwd()
			if err != nil {
				http.Error(w, "Error getting current working directory", http.StatusInternalServerError)
				return
			}
			filePath := filepath.Join(cwd, "output", a.Series.Suffix, "annotated", address+".png")
			if _, err := os.Stat(filePath); os.IsNotExist(err) {
				msg := fmt.Sprintf("File not found at %s", filePath)
				http.Error(w, msg, http.StatusNotFound)
				return
			}
			w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
			w.Header().Set("Pragma", "no-cache")
			w.Header().Set("Expires", "0")
			http.ServeFile(w, r, filePath)
		})

		go func() {
			logger.Info("Starting file server on :8889")
			if err := http.ListenAndServe(":8889", nil); err != nil {
				logger.Error("File server error:", err)
			}
		}()

		go a.StartServers()

		if err := wails.Run(&opts); err != nil {
			fmt.Println("Error:", err.Error())
		}
	}
}
