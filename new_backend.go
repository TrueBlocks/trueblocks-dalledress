package main

import (
	"fmt"

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
		if dd.EnhancedPrompt, err = openai.EnhancePrompt(app.GetPrompt(addr)); err != nil {
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

func (app *App) GetImage(addr string) {
	if dd, err := NewDalleDress(app.databases, addr); err != nil {
		logger.Error(err.Error())
	} else {
		imageData := openai.ImageData{
			EnhancedPrompt: app.GetEnhanced(addr),
			TersePrompt:    app.GetTerse(addr),
			Hash:           dd.Orig,
			SeriesName:     app.Series.Suffix,
		}
		if err := openai.GetImage(&imageData); err != nil {
			logger.Error(err.Error())
		}
		fmt.Println("Open address: ", addr)
	}
}
