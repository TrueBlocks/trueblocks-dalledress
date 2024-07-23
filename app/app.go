package app

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand/v2"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"text/template"
	"time"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/base"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/file"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/logger"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/output"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/config"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/dalle"
	"github.com/joho/godotenv"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type App struct {
	ctx        context.Context
	session    config.Session
	apiKeys    map[string]string
	namesMap   map[base.Address]NameEx
	names      []NameEx // We keep both for performance reasons
	ensMap     map[string]base.Address
	renderCtxs map[base.Address][]*output.RenderCtx
	// Add your application's data here
	databases      map[string][]string
	authorTemplate *template.Template
	promptTemplate *template.Template
	dataTemplate   *template.Template
	terseTemplate  *template.Template
	titleTemplate  *template.Template
	Series         dalle.Series `json:"series"`
	dalleCache     map[string]*dalle.DalleDress
}

func NewApp() *App {
	a := App{
		apiKeys:    make(map[string]string),
		namesMap:   make(map[base.Address]NameEx),
		renderCtxs: make(map[base.Address][]*output.RenderCtx),
		ensMap:     make(map[string]base.Address),
		// Initialize maps here
		databases:  make(map[string][]string),
		dalleCache: make(map[string]*dalle.DalleDress),
	}

	// it's okay if it's not found
	_ = a.session.Load()

	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	} else if a.apiKeys["openAi"] = os.Getenv("OPENAI_API_KEY"); a.apiKeys["openAi"] == "" {
		log.Fatal("No OPENAI_API_KEY key found")
	}

	// Initialize your data here
	var err error
	if a.promptTemplate, err = template.New("prompt").Parse(promptTemplate); err != nil {
		logger.Fatal("could not create prompt template:", err)
	}
	if a.dataTemplate, err = template.New("data").Parse(dataTemplate); err != nil {
		logger.Fatal("could not create data template:", err)
	}
	if a.titleTemplate, err = template.New("terse").Parse(titleTemplate); err != nil {
		logger.Fatal("could not create title template:", err)
	}
	if a.terseTemplate, err = template.New("terse").Parse(terseTemplate); err != nil {
		logger.Fatal("could not create terse template:", err)
	}
	if a.authorTemplate, err = template.New("author").Parse(authorTemplate); err != nil {
		logger.Fatal("could not create prompt template:", err)
	}
	logger.Info()
	logger.Info("Compiled templates")

	a.ReloadDatabases()

	return &a
}

func (a App) String() string {
	bytes, _ := json.MarshalIndent(a, "", "  ")
	return string(bytes)
}

func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx
	if err := a.loadNames(); err != nil {
		logger.Panic(err)
	}
}

func (a *App) DomReady(ctx context.Context) {
	// Sometimes useful for debugging
	if os.Getenv("TB_CMD_LINE") == "true" {
		return
	}
	runtime.WindowSetPosition(a.ctx, a.session.X, a.session.Y)
	runtime.WindowSetSize(a.ctx, a.session.Width, a.session.Height)
	runtime.WindowShow(a.ctx)
}

func (a *App) Shutdown(ctx context.Context) {
	// Sometimes useful for debugging
	if os.Getenv("TB_CMD_LINE") == "true" {
		return
	}
	a.session.X, a.session.Y = runtime.WindowGetPosition(a.ctx)
	a.session.Width, a.session.Height = runtime.WindowGetSize(a.ctx)
	a.session.Y += 38 // TODO: This is a hack to account for the menu bar - not sure why it's needed
	a.session.Save()
}

func (a *App) GetSession() *config.Session {
	return &a.session
}

func (a *App) ReloadDatabases() {
	a.Series = dalle.Series{}
	a.databases = make(map[string][]string)

	var err error
	if a.Series, err = a.LoadSeries(); err != nil {
		logger.Fatal(err)
	}
	logger.Info("Loaded series:", a.Series.Suffix)

	for _, db := range dalle.DatabaseNames {
		if a.databases[db] == nil {
			if lines, err := a.toLines(db); err != nil {
				logger.Fatal(err)
			} else {
				a.databases[db] = lines
				for i := 0; i < len(a.databases[db]); i++ {
					a.databases[db][i] = strings.Replace(a.databases[db][i], "v0.1.0,", "", -1)
				}
			}
		}
	}
	logger.Info("Loaded", len(dalle.DatabaseNames), "databases")
}

func (a *App) LoadSeries() (dalle.Series, error) {
	lastSeries := a.GetSession().LastSeries
	fn := filepath.Join("./output/series", lastSeries+".json")
	str := strings.TrimSpace(file.AsciiFileToString(fn))
	logger.Info("lastSeries", lastSeries)
	if len(str) == 0 || !file.FileExists(fn) {
		logger.Info("No series found, creating a new one", fn)
		ret := dalle.Series{
			Suffix: "simple",
		}
		ret.SaveSeries(fn, 0)
		return ret, nil
	}

	bytes := []byte(str)
	var s dalle.Series
	if err := json.Unmarshal(bytes, &s); err != nil {
		logger.Error("could not unmarshal series:", err)
		return dalle.Series{}, err
	}

	s.Suffix = strings.Trim(strings.ReplaceAll(s.Suffix, " ", "-"), "-")
	s.SaveSeries(filepath.Join("./output/series", s.Suffix+".json"), 0)
	return s, nil
}

func (a *App) toLines(db string) ([]string, error) {
	filename := "./databases/" + db + ".csv"
	lines := file.AsciiFileToLines(filename)
	lines = lines[1:] // skip header
	var err error
	if len(lines) == 0 {
		err = fmt.Errorf("could not load %s", filename)
	} else {
		fn := strings.ToUpper(db[:1]) + db[1:]
		if filter, err := a.Series.GetFilter(fn); err != nil {
			return lines, err

		} else {
			if len(filter) == 0 {
				return lines, nil
			}

			filtered := make([]string, 0, len(lines))
			for _, line := range lines {
				for _, f := range filter {
					if strings.Contains(line, f) {
						filtered = append(filtered, line)
					}
				}
			}
			lines = filtered
		}
	}

	if len(lines) == 0 {
		lines = append(lines, "none")
	}

	return lines, err
}

func (a *App) HandleLines() {
	batchSize := 5
	rateLimit := time.Second / 5
	sem := make(chan struct{}, batchSize)

	lines, series := func() ([]string, []dalle.Series) {
		series := []dalle.Series{a.Series}
		if len(os.Args) < 2 {
			return file.AsciiFileToLines("inputs/addresses.txt"), series
		}
		return os.Args[1:], series
	}()

	var wg sync.WaitGroup
	ticker := time.NewTicker(rateLimit)
	defer ticker.Stop()

	for _, ser := range series {
		a.Series = ser
		for i, addr := range lines {
			if a.Series.Last > 0 && i <= int(a.Series.Last) {
				continue
			}

			sem <- struct{}{}
			wg.Add(1)
			go func(index int, address string) {
				defer wg.Done()
				defer func() { <-sem }()
				backoff := time.Second
				maxRetries := 5
				for attempt := 0; attempt < maxRetries; attempt++ {
					<-ticker.C
					_, err := a.GenerateImage(address)
					if err == nil {
						return
					}
					// if isContentPolicyViolation(err) {
					// 	msg := fmt.Sprintf("Content policy violation, skipping retry for address: %s Error: %s", address, err)
					// 	logger.Error(msg)
					// 	return
					// } else
					if strings.Contains(err.Error(), "seed length is less than 66") {
						msg := fmt.Sprintf("Invalid address, skipping retry for address: %s Error: %s", address, err)
						logger.Error(msg)
						return
					}
					msg := fmt.Sprintf("Error fetching image: %s Retry attempt: %d Sleeping: %d", err, attempt+1, backoff)
					logger.Error(msg)
					time.Sleep(backoff)
					backoff = time.Duration(float64(backoff) * (1 + rand.Float64()))
				}
				logger.Error("Failed to fetch image after max retries:", address)
			}(i, addr)

			a.Series.Last = i
			a.Series.SaveSeries("inputs/series.json", a.Series.Last)

			if (i+1)%batchSize == 0 {
				wg.Wait()
			}
		}
	}
	wg.Wait()
}
