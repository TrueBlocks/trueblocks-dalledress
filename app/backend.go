package app

import (
	"path/filepath"
	"strings"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/logger"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/dalle"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/openai"
)

func (a *App) GetPrompt(addr string) string {
	if dd, err := dalle.NewDalleDress(a.databases, addr); err != nil {
		return err.Error()
	} else {
		dd.Prompt, _ = dd.GeneratePrompt(a.pTemplate, nil)
		return dd.Prompt
	}
}

func (a *App) GetEnhanced(addr string) string {
	if dd, err := dalle.NewDalleDress(a.databases, addr); err != nil {
		return err.Error()
	} else {
		authorType, _ := dd.GeneratePrompt(a.aTemplate, nil)
		if dd.EnhancedPrompt, err = openai.EnhancePrompt(a.GetPrompt(addr), authorType); err != nil {
			logger.Fatal(err.Error())
		}
		msg := " DO NOT PUT TEXT IN THE IMAGE. "
		dd.EnhancedPrompt = msg + dd.EnhancedPrompt + msg
		return dd.EnhancedPrompt
	}
}

func (a *App) GetTerse(addr string) string {
	if dd, err := dalle.NewDalleDress(a.databases, addr); err != nil {
		return err.Error()
	} else {
		dd.TersePrompt, _ = dd.GeneratePrompt(a.tTemplate, nil)
		return dd.TersePrompt
	}
}

func (a *App) GetData(addr string) string {
	if dd, err := dalle.NewDalleDress(a.databases, addr); err != nil {
		return err.Error()
	} else {
		dd.DataPrompt, _ = dd.GeneratePrompt(a.dTemplate, nil)
		return dd.DataPrompt
	}
}

func (a *App) GetJson(addr string) string {
	if dd, err := dalle.NewDalleDress(a.databases, addr); err != nil {
		return err.Error()
	} else {
		return dd.String()
	}
}

func (a *App) GetImage(addr string) error {
	if dd, err := dalle.NewDalleDress(a.databases, addr); err != nil {
		return err
	} else {
		suff := a.Series.Suffix
		ReportOn(filepath.Join(suff, "selector"), dd.Orig, "json", dd.String())
		data := a.GetData(addr)
		ReportOn(filepath.Join(suff, "data"), dd.Orig, "txt", data)
		terse := a.GetTerse(addr)
		ReportOn(filepath.Join(suff, "terse"), dd.Orig, "txt", terse)
		prompt := a.GetPrompt(addr)
		ReportOn(filepath.Join(suff, "prompt"), dd.Orig, "txt", prompt)
		enhanced := a.GetEnhanced(addr)
		ReportOn(filepath.Join(suff, "enhanced"), dd.Orig, "txt", strings.ReplaceAll(strings.ReplaceAll(enhanced, ".", ".\n"), ",", ",\n"))
		imageData := openai.ImageData{
			TersePrompt:    terse,
			EnhancedPrompt: enhanced,
			SeriesName:     a.Series.Suffix,
			Filename:       strings.ReplaceAll(dd.Orig, ",", "-"),
		}
		return openai.GetImage(&imageData)
	}
}
