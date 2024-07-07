/*
GetJson,
GetData,
GetSeries,
GetTerse,
GetPrompt,
GetEnhanced,
GetImage,
GenerateImage,
Refresh,
*/

package app

import (
	"path/filepath"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/file"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/logger"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/dalle"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/openai"
)

func (a *App) GetSeries(addr string) string {
	return a.Series.String()
}

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
		fn := filepath.Join("output", a.Series.Suffix, "enhanced", dd.Filename+".txt")
		if file.FileExists(fn) {
			return file.AsciiFileToString(fn)
		} else {
			return "No enhanced prompt found at " + fn + ". Press Generate."
		}
	}
}

func (a *App) GetImage(addr string) (string, error) {
	if dd, err := dalle.NewDalleDress(a.databases, addr); err != nil {
		return err.Error(), err
	} else {
		return dd.Filename, nil
	}
}

func (a *App) GenerateEnhanced(addr string) string {
	if dd, err := dalle.NewDalleDress(a.databases, addr); err != nil {
		return err.Error()
	} else {
		authorType, _ := dd.ExecuteTemplate(a.authorTemplate, nil)
		if dd.EnhancedPrompt, err = openai.EnhancePrompt(a.GetPrompt(addr), authorType); err != nil {
			logger.Fatal(err.Error())
		}
		msg := " DO NOT PUT TEXT IN THE IMAGE. "
		dd.EnhancedPrompt = msg + dd.EnhancedPrompt + msg
		return dd.EnhancedPrompt
	}
}

func (a *App) Refresh(addr string) string {
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
		return ""
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
		enhanced := a.GenerateEnhanced(addr)
		dd.ReportOn(filepath.Join(suff, "enhanced"), "txt", enhanced)
		imageData := openai.ImageData{
			TitlePrompt:    title,
			TersePrompt:    terse,
			EnhancedPrompt: enhanced,
			SeriesName:     a.Series.Suffix,
			Filename:       dd.Filename,
		}
		if err := openai.GenerateImage(&imageData); err != nil {
			return err.Error()
		}
		return enhanced
	}
}

func (a *App) GetExistingAddrs() []string {
	return []string{
		"gitcoin.eth",
		"giveth.eth",
		"chase.wright.eth",
		"cnn.eth",
		"dawid.eth",
		"dragonstone.eth",
		"eats.eth",
		"ens.eth",
		"gameofthrones.eth",
		"jen.eth",
		"makingprogress.eth",
		"meriam.eth",
		"nate.eth",
		"poap.eth",
		"revenge.eth",
		"rotki.eth",
		"trueblocks.eth",
		"unchainedindex.eth",
		"vitalik.eth",
		"when.eth",
	}
}
