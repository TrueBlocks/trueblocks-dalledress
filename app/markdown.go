package app

import (
	"path/filepath"

	"github.com/TrueBlocks/trueblocks-dalledress/pkg/markdown"
)

func (a *App) GetMarkdown(folder, route, tab string) string {
	lang := a.Preferences.App.Language
	if md, err := markdown.LoadMarkdown(a.Assets, filepath.Join("frontend", "src", "assets", folder), lang, route, tab); err != nil {
		return err.Error()
	} else {
		return md
	}
}
