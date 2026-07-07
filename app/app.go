// Package app holds the Wails-bound application object. The embedded NavState
// provides the standard persistence surface (GetTab/SetTab, GetLastRoute/
// SetLastRoute, window geometry, sidebar width, prefs, table states) as bound
// methods; see ai/Architecture.md §7 in the trueblocks-art repo.
package app

import (
	"context"

	appkit "github.com/TrueBlocks/trueblocks-art/packages/appkit/v2"
	dalle "github.com/TrueBlocks/trueblocks-dalle/v6"
	"github.com/TrueBlocks/trueblocks-dalle/v6/pkg/storage"
)

type App struct {
	*appkit.NavState
	ctx    context.Context
	engine *dalle.Engine
}

func NewApp(prefsPath string) (*App, error) {
	engine, err := dalle.New(dalle.Config{})
	if err != nil {
		return nil, err
	}
	return &App{
		NavState: appkit.NewNavState(prefsPath, appkit.NavDefaults{
			Route:        "dashboard",
			SidebarWidth: 220,
		}),
		engine: engine,
	}, nil
}

func (a *App) Startup(ctx context.Context) { a.ctx = ctx }

func (a *App) Shutdown(_ context.Context) {}

func (a *App) Preview(request dalle.GenerateRequest) (dalle.GenerateResult, error) {
	return a.engine.Preview(request)
}

func (a *App) Generate(request dalle.GenerateRequest) (dalle.GenerateResult, error) {
	return a.engine.Generate(request)
}

func (a *App) ListImages(series string) ([]dalle.ImageMetadataRecord, error) {
	return a.engine.ListImages(dalle.ImageFilter{Series: series})
}

func (a *App) GetImage(id string) (dalle.ImageMetadataRecord, error) {
	return a.engine.GetImage(id)
}

func (a *App) ExportImage(id string, options dalle.ExportImageOptions) (dalle.ExportImageResult, error) {
	return a.engine.ExportImage(id, options)
}

func (a *App) ListSeries(includeHidden bool, onlyHidden bool) ([]dalle.Series, error) {
	return a.engine.ListSeries(dalle.SeriesFilter{IncludeHidden: includeHidden, OnlyHidden: onlyHidden})
}

func (a *App) GetSeries(name string) (dalle.Series, error) {
	return a.engine.GetSeries(name)
}

func (a *App) SaveSeries(series dalle.Series) (dalle.Series, error) {
	return a.engine.SaveSeries(series)
}

func (a *App) SetSeriesHidden(name string, hidden bool) (dalle.Series, error) {
	return a.engine.SetSeriesHidden(name, hidden)
}

func (a *App) ListDatabaseArchives() ([]storage.DatabaseArchiveManifest, error) {
	return a.engine.ListDatabaseArchives()
}

func (a *App) GetDatabaseArchive(version string) (storage.DatabaseArchiveManifest, error) {
	return a.engine.GetDatabaseArchive(version)
}

func (a *App) ValidateDalle() error {
	return a.engine.Validate()
}
