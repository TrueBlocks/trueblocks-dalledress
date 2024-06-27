package main

import (
	"path/filepath"
	"strings"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/logger"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/openai"
)

func (app *App) GetPrompt(addr string) string {
	if dd, err := NewDalleDress(app.databases, addr); err != nil {
		return err.Error()
	} else {
		dd.Prompt, _ = dd.generatePrompt(app.pTemplate, nil)
		return dd.Prompt
	}
}

func (app *App) GetEnhanced(addr string) string {
	if dd, err := NewDalleDress(app.databases, addr); err != nil {
		return err.Error()
	} else {
		authorType, _ := dd.generatePrompt(app.aTemplate, nil)
		if dd.EnhancedPrompt, err = openai.EnhancePrompt(app.GetPrompt(addr), authorType); err != nil {
			logger.Fatal(err.Error())
		}
		msg := " DO NOT PUT TEXT IN THE IMAGE. "
		dd.EnhancedPrompt = msg + dd.EnhancedPrompt + msg
		return dd.EnhancedPrompt
	}
}

func (app *App) GetTerse(addr string) string {
	if dd, err := NewDalleDress(app.databases, addr); err != nil {
		return err.Error()
	} else {
		dd.TersePrompt, _ = dd.generatePrompt(app.tTemplate, nil)
		return dd.TersePrompt
	}
}

func (app *App) GetData(addr string) string {
	if dd, err := NewDalleDress(app.databases, addr); err != nil {
		return err.Error()
	} else {
		dd.DataPrompt, _ = dd.generatePrompt(app.dTemplate, nil)
		return dd.DataPrompt
	}
}

func (app *App) GetJson(addr string) string {
	if dd, err := NewDalleDress(app.databases, addr); err != nil {
		return err.Error()
	} else {
		return dd.String()
	}
}

func (app *App) GetImage(addr string) error {
	if dd, err := NewDalleDress(app.databases, addr); err != nil {
		return err
	} else {
		suff := app.Series.Suffix
		ReportOn(filepath.Join(suff, "selector"), dd.Orig, "json", dd.String())
		data := app.GetData(addr)
		ReportOn(filepath.Join(suff, "data"), dd.Orig, "txt", data)
		terse := app.GetTerse(addr)
		ReportOn(filepath.Join(suff, "terse"), dd.Orig, "txt", terse)
		prompt := app.GetPrompt(addr)
		ReportOn(filepath.Join(suff, "prompt"), dd.Orig, "txt", prompt)
		enhanced := app.GetEnhanced(addr)
		ReportOn(filepath.Join(suff, "enhanced"), dd.Orig, "txt", strings.ReplaceAll(enhanced, ".", ".\n"))
		imageData := openai.ImageData{
			TersePrompt:    terse,
			EnhancedPrompt: enhanced,
			SeriesName:     app.Series.Suffix,
			Filename:       strings.ReplaceAll(dd.Orig, ",", "-"),
		}
		return openai.GetImage(&imageData)
	}
}
