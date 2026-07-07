// Package app holds the Wails-bound application object. The embedded NavState
// provides the standard persistence surface (GetTab/SetTab, GetLastRoute/
// SetLastRoute, window geometry, sidebar width, prefs, table states) as bound
// methods; see ai/Architecture.md §7 in the trueblocks-art repo.
package app

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"os/exec"

	appkit "github.com/TrueBlocks/trueblocks-art/packages/appkit/v2"
	dalle "github.com/TrueBlocks/trueblocks-dalle/v6"
	"github.com/TrueBlocks/trueblocks-dalle/v6/pkg/progress"
	"github.com/TrueBlocks/trueblocks-dalle/v6/pkg/storage"
)

type App struct {
	*appkit.NavState
	ctx    context.Context
	engine *dalle.Engine
}

type GenerationProgress struct {
	Active     bool    `json:"active"`
	Series     string  `json:"series"`
	Seed       string  `json:"seed"`
	Phase      string  `json:"phase"`
	Percent    float64 `json:"percent"`
	ETASeconds float64 `json:"etaSeconds"`
	Done       bool    `json:"done"`
	CacheHit   bool    `json:"cacheHit"`
	Error      string  `json:"error"`
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

func (a *App) GetGenerationProgress(series string, seed string) GenerationProgress {
	report := progress.GetProgress(series, seed)
	if report == nil {
		return GenerationProgress{Series: series, Seed: seed}
	}
	return GenerationProgress{
		Active:     true,
		Series:     report.Series,
		Seed:       report.Address,
		Phase:      string(report.Current),
		Percent:    report.Percent,
		ETASeconds: report.ETASeconds,
		Done:       report.Done,
		CacheHit:   report.CacheHit,
		Error:      report.Error,
	}
}

func (a *App) ListImages(series string) ([]dalle.ImageMetadataRecord, error) {
	return a.engine.ListImages(dalle.ImageFilter{Series: series})
}

func (a *App) GetImage(id string) (dalle.ImageMetadataRecord, error) {
	return a.engine.GetImage(id)
}

func (a *App) GetImageArtifactDataURL(id string, artifact string) (string, error) {
	path, err := a.imageArtifactPath(id, artifact)
	if err != nil {
		return "", err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return "data:image/png;base64," + base64.StdEncoding.EncodeToString(data), nil
}

func (a *App) OpenImageArtifact(id string, artifact string) error {
	path, err := a.imageArtifactPath(id, artifact)
	if err != nil {
		return err
	}
	return exec.Command("open", path).Run()
}

func (a *App) RevealImageArtifact(id string, artifact string) error {
	path, err := a.imageArtifactPath(id, artifact)
	if err != nil {
		return err
	}
	return exec.Command("open", "-R", path).Run()
}

func (a *App) imageArtifactPath(id string, artifact string) (string, error) {
	record, err := a.engine.GetImage(id)
	if err != nil {
		return "", err
	}
	path := ""
	switch artifact {
	case "annotated":
		path = record.Metadata.Artifacts.Annotated
	case "generated":
		path = record.Metadata.Artifacts.Generated
	default:
		return "", fmt.Errorf("unknown artifact %q", artifact)
	}
	if path == "" {
		return "", fmt.Errorf("%s artifact is not available", artifact)
	}
	return path, nil
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
