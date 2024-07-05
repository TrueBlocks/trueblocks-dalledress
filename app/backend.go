package app

import (
	"encoding/base64"
	"os"
	"path/filepath"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/file"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/logger"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/dalle"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/openai"
)

func (a *App) GetJson(addr string) string {
	if dd, err := dalle.NewDalleDress(a.databases, addr); err != nil {
		return err.Error()
	} else {
		return dd.String()
	}
}

func (a *App) GetData(addr string) string {
	if dd, err := dalle.NewDalleDress(a.databases, addr); err != nil {
		return err.Error()
	} else {
		dd.DataPrompt, _ = dd.ExecuteTemplate(a.dataTemplate, nil)
		return dd.DataPrompt
	}
}

func (a *App) GetTitle(addr string) string {
	if dd, err := dalle.NewDalleDress(a.databases, addr); err != nil {
		return err.Error()
	} else {
		dd.TitlePrompt, _ = dd.ExecuteTemplate(a.titleTemplate, nil)
		return dd.TitlePrompt
	}
}

func (a *App) GetTerse(addr string) string {
	if dd, err := dalle.NewDalleDress(a.databases, addr); err != nil {
		return err.Error()
	} else {
		dd.TersePrompt, _ = dd.ExecuteTemplate(a.terseTemplate, nil)
		return dd.TersePrompt
	}
}

func (a *App) GetPrompt(addr string) string {
	if dd, err := dalle.NewDalleDress(a.databases, addr); err != nil {
		return err.Error()
	} else {
		dd.Prompt, _ = dd.ExecuteTemplate(a.promptTemplate, nil)
		return dd.Prompt
	}
}

func (a *App) GetEnhanced(addr string) string {
	if dd, err := dalle.NewDalleDress(a.databases, addr); err != nil {
		return err.Error()
	} else {
		fn := filepath.Join("output", a.Series.Suffix, "enhanced", dd.FileName+".txt")
		logger.Info("Loading file", fn)
		if file.FileExists(fn) {
			return file.AsciiFileToString(fn)
		} else {
			return "No enhanced prompt found at " + fn + ". Press Generate."
		}
	}
}

func (a *App) GetImage(addr string) (string, error) {
	logger.Info("GetImage", addr)
	if dd, err := dalle.NewDalleDress(a.databases, addr); err != nil {
		return err.Error(), err
	} else {
		cwd, err := os.Getwd()
		if err != nil {
			return err.Error(), err
		}
		fn := filepath.Join(cwd, "output", a.Series.Suffix, "annotated", dd.FileName+".png")
		logger.Info("Loading file", fn)
		if file.FileExists(fn) {
			logger.Info("loading image from:", fn)
			imageBytes, err := os.ReadFile(fn)
			if err != nil {
				return err.Error(), err
			}
			base64Image := base64.StdEncoding.EncodeToString(imageBytes)
			return base64Image, nil
		} else {
			logger.Info("File not found", fn)
			return "No image file found at" + fn + ". Press Generate.", nil
		}
	}
}

func (a *App) GenerateImage(addr string) string {
	if dd, err := dalle.NewDalleDress(a.databases, addr); err != nil {
		return err.Error()
	} else {
		suff := a.Series.Suffix
		dd.ReportOn(filepath.Join(suff, "selector"), "json", dd.String())
		data := a.GetData(addr)
		dd.ReportOn(filepath.Join(suff, "data"), "txt", data)
		title := a.GetTitle(addr)
		dd.ReportOn(filepath.Join(suff, "title"), "txt", title)
		terse := a.GetTerse(addr)
		dd.ReportOn(filepath.Join(suff, "terse"), "txt", terse)
		prompt := a.GetPrompt(addr)
		dd.ReportOn(filepath.Join(suff, "prompt"), "txt", prompt)
		enhanced := a.GetEnhanced(addr)
		dd.ReportOn(filepath.Join(suff, "enhanced"), "txt", enhanced)
		imageData := openai.ImageData{
			TitlePrompt:    title,
			TersePrompt:    terse,
			EnhancedPrompt: enhanced,
			SeriesName:     a.Series.Suffix,
			Filename:       dd.FileName,
		}
		if err := openai.GenerateImage(&imageData); err != nil {
			return err.Error()
		}
		return enhanced
	}
}
