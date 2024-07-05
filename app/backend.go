package app

import (
	"encoding/base64"
	"fmt"
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
		dd.DataPrompt, _ = dd.GeneratePrompt(a.dataTemplate, nil)
		return dd.DataPrompt
	}
}

func (a *App) GetTitle(addr string) string {
	if dd, err := dalle.NewDalleDress(a.databases, addr); err != nil {
		return err.Error()
	} else {
		dd.TitlePrompt, _ = dd.GeneratePrompt(a.titleTemplate, nil)
		return dd.TitlePrompt
	}
}

func (a *App) GetTerse(addr string) string {
	if dd, err := dalle.NewDalleDress(a.databases, addr); err != nil {
		return err.Error()
	} else {
		dd.TersePrompt, _ = dd.GeneratePrompt(a.terseTemplate, nil)
		return dd.TersePrompt
	}
}

func (a *App) GetPrompt(addr string) string {
	if dd, err := dalle.NewDalleDress(a.databases, addr); err != nil {
		return err.Error()
	} else {
		dd.Prompt, _ = dd.GeneratePrompt(a.promptTemplate, nil)
		return dd.Prompt
	}
}

func (a *App) GetEnhanced(addr string) string {
	if dd, err := dalle.NewDalleDress(a.databases, addr); err != nil {
		return err.Error()
	} else {
		cwd, err := os.Getwd()
		if err != nil {
			return err.Error()
		}
		fn := filepath.Join(cwd, "output", a.Series.Suffix, "enhanced", dd.Orig+".txt")
		if file.FileExists(fn) {
			return fn + file.AsciiFileToString(fn)
		} else {
			return "No enhanced prompt found. Press Generate."
		}
	}
}

func (a *App) GetImage(addr string) (string, error) {
	if dd, err := dalle.NewDalleDress(a.databases, addr); err != nil {
		return err.Error(), err
	} else {
		cwd, err := os.Getwd()
		if err != nil {
			return err.Error(), err
		}
		fn := fmt.Sprintf("%s/output/%s/annotated/%s.png", cwd, a.Series.Suffix, dd.Orig)
		if file.FileExists(fn) {
			logger.Info("loading image from:", fn)
			imageBytes, err := os.ReadFile(fn)
			if err != nil {
				return err.Error(), err
			}
			base64Image := base64.StdEncoding.EncodeToString(imageBytes)
			return base64Image, nil
		} else {
			return "No image file found. Press Generate.", nil
		}
	}
}

func (a *App) GenerateEnhanced(addr string) string {
	if dd, err := dalle.NewDalleDress(a.databases, addr); err != nil {
		return err.Error()
	} else {
		authorType, _ := dd.GeneratePrompt(a.authorTemplate, nil)
		if dd.EnhancedPrompt, err = openai.EnhancePrompt(a.GetPrompt(addr), authorType); err != nil {
			logger.Fatal(err.Error())
		}
		msg := " DO NOT PUT TEXT IN THE IMAGE. "
		dd.EnhancedPrompt = msg + dd.EnhancedPrompt + msg
		return dd.EnhancedPrompt
	}
}

/*
	if dd, err := dalle.NewDalleDress(a.databases, addr); err != nil {
		return err
	} else {
		suff := a.Series.Suffix
		ReportOn(filepath.Join(suff, "selector"), dd.Orig, "json", dd.String())
		data := a.GetData(addr)
		ReportOn(filepath.Join(suff, "data"), dd.Orig, "txt", data)
		title := a.GetTitle(addr)
		ReportOn(filepath.Join(suff, "title"), dd.Orig, "txt", title)
		terse := a.GetTerse(addr)
		ReportOn(filepath.Join(suff, "terse"), dd.Orig, "txt", terse)
		prompt := a.GetPrompt(addr)
		ReportOn(filepath.Join(suff, "prompt"), dd.Orig, "txt", prompt)
		enhanced := a.GetEnhanced(addr)
		ReportOn(filepath.Join(suff, "enhanced"), dd.Orig, "txt", strings.ReplaceAll(strings.ReplaceAll(enhanced, ".", ".\n"), ",", ",\n"))
		imageData := openai.ImageData{
			TitlePrompt:    title,
			TersePrompt:    terse,
			EnhancedPrompt: enhanced,
			SeriesName:     a.Series.Suffix,
			Filename:       strings.ReplaceAll(dd.Orig, ",", "-"),
		}
		return openai.GetImage(&imageData)
	}
*/
