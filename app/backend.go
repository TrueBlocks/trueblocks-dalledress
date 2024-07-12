package app

import (
	"fmt"
	"path/filepath"
	"strings"
	"sync"

	"github.com/TrueBlocks/trueblocks-core/sdk"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/file"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/logger"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/dalle"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/openai"
)

var dalleCacheMutex sync.Mutex

func (a *App) MakeDalleDress(addressIn string) (*dalle.DalleDress, error) {
	dalleCacheMutex.Lock()
	defer dalleCacheMutex.Unlock()
	if a.dalleCache[addressIn] != nil {
		logger.Info("Returning cached dalle for", addressIn)
		return a.dalleCache[addressIn], nil
	}

	address := addressIn
	logger.Info("Making dalle for", addressIn)
	if strings.HasSuffix(address, ".eth") {
		opts := sdk.NamesOptions{
			Terms: []string{address},
		}
		if names, _, err := opts.Names(); err != nil {
			return nil, fmt.Errorf("error getting names for %s", address)
		} else {
			if len(names) > 0 {
				address = names[0].Address.Hex()
			}
		}
	}
	logger.Info("Resolved", addressIn)

	parts := strings.Split(address, ",")
	seed := parts[0] + reverse(parts[0])
	if len(seed) < 66 {
		return nil, fmt.Errorf("seed length is less than 66")
	}
	if strings.HasPrefix(seed, "0x") {
		seed = seed[2:66]
	}

	fn := validFilename(address)
	if a.dalleCache[fn] != nil {
		logger.Info("Returning cached dalle for", addressIn)
		return a.dalleCache[fn], nil
	}

	dd := dalle.DalleDress{
		Original:  addressIn,
		Filename:  fn,
		Seed:      seed,
		AttribMap: make(map[string]dalle.Attribute),
	}

	for i := 0; i < len(dd.Seed); i = i + 8 {
		index := len(dd.Attribs)
		attr := dalle.NewAttribute(a.databases, index, dd.Seed[i:i+6])
		dd.Attribs = append(dd.Attribs, attr)
		dd.AttribMap[attr.Name] = attr
		if i+4+6 < len(dd.Seed) {
			index = len(dd.Attribs)
			attr = dalle.NewAttribute(a.databases, index, dd.Seed[i+4:i+4+6])
			dd.Attribs = append(dd.Attribs, attr)
			dd.AttribMap[attr.Name] = attr
		}
	}

	suff := a.Series.Suffix
	dd.DataPrompt, _ = dd.ExecuteTemplate(a.dataTemplate, nil)
	dd.ReportOn(addressIn, filepath.Join(suff, "data"), "txt", dd.DataPrompt)
	dd.TitlePrompt, _ = dd.ExecuteTemplate(a.titleTemplate, nil)
	dd.ReportOn(addressIn, filepath.Join(suff, "title"), "txt", dd.TitlePrompt)
	dd.TersePrompt, _ = dd.ExecuteTemplate(a.terseTemplate, nil)
	dd.ReportOn(addressIn, filepath.Join(suff, "terse"), "txt", dd.TersePrompt)
	dd.Prompt, _ = dd.ExecuteTemplate(a.promptTemplate, nil)
	dd.ReportOn(addressIn, filepath.Join(suff, "prompt"), "txt", dd.Prompt)
	fn = filepath.Join("output", a.Series.Suffix, "enhanced", dd.Filename+".txt")
	dd.EnhancedPrompt = ""
	if file.FileExists(fn) {
		dd.EnhancedPrompt = file.AsciiFileToString(fn)
	}

	a.dalleCache[dd.Filename] = &dd
	a.dalleCache[addressIn] = &dd

	return &dd, nil
}

func (a *App) GetSeries(addr string) string {
	return a.Series.String()
}

func (a *App) GetJson(addr string) string {
	if dd, err := a.MakeDalleDress(addr); err != nil {
		return err.Error()
	} else {
		return dd.String()
	}
}

func (a *App) GetData(addr string) string {
	if dd, err := a.MakeDalleDress(addr); err != nil {
		return err.Error()
	} else {
		return dd.DataPrompt
	}
}

func (a *App) GetTitle(addr string) string {
	if dd, err := a.MakeDalleDress(addr); err != nil {
		return err.Error()
	} else {
		return dd.TitlePrompt
	}
}

func (a *App) GetTerse(addr string) string {
	if dd, err := a.MakeDalleDress(addr); err != nil {
		return err.Error()
	} else {
		return dd.TersePrompt
	}
}

func (a *App) GetPrompt(addr string) string {
	if dd, err := a.MakeDalleDress(addr); err != nil {
		return err.Error()
	} else {
		return dd.Prompt
	}
}

func (a *App) GetEnhanced(addr string) string {
	if dd, err := a.MakeDalleDress(addr); err != nil {
		return err.Error()
	} else {
		return dd.EnhancedPrompt
	}
}

func (a *App) GetFilename(addr string) string {
	if dd, err := a.MakeDalleDress(addr); err != nil {
		return err.Error()
	} else {
		return dd.Filename
	}
}

func (a *App) Save(addr string) bool {
	if dd, err := a.MakeDalleDress(addr); err != nil {
		return false
	} else {
		dd.ReportOn(addr, filepath.Join(a.Series.Suffix, "selector"), "json", dd.String())
		return true
	}
}

func (a *App) GenerateEnhanced(addr string) string {
	if dd, err := a.MakeDalleDress(addr); err != nil {
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

func (a *App) GenerateImage(addr string) (string, error) {
	if dd, err := a.MakeDalleDress(addr); err != nil {
		return err.Error(), err
	} else {
		suff := a.Series.Suffix
		dd.EnhancedPrompt = a.GenerateEnhanced(addr)
		dd.ReportOn(addr, filepath.Join(suff, "enhanced"), "txt", dd.EnhancedPrompt)
		_ = a.Save(addr)
		imageData := openai.ImageData{
			TitlePrompt:    dd.TitlePrompt,
			TersePrompt:    dd.TersePrompt,
			EnhancedPrompt: dd.EnhancedPrompt,
			SeriesName:     a.Series.Suffix,
			Filename:       dd.Filename,
		}
		if err := openai.RequestImage(&imageData); err != nil {
			return err.Error(), err
		}
		return dd.EnhancedPrompt, nil
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

// validFilename returns a valid filename from the input string
func validFilename(in string) string {
	invalidChars := []string{"/", "\\", ":", "*", "?", "\"", "<", ">", "|"}
	for _, char := range invalidChars {
		in = strings.ReplaceAll(in, char, "_")
	}
	in = strings.TrimSpace(in)
	in = strings.ReplaceAll(in, "__", "_")
	return in
}

// reverse returns the reverse of the input string
func reverse(s string) string {
	runes := []rune(s)
	n := len(runes)
	for i := 0; i < n/2; i++ {
		runes[i], runes[n-1-i] = runes[n-1-i], runes[i]
	}
	return string(runes)
}
