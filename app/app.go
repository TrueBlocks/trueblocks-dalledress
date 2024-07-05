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

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/file"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/logger"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/types"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/config"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/dalle"
	"github.com/joho/godotenv"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type App struct {
	ctx            context.Context
	session        config.Session
	apiKey         string
	databases      map[string][]string
	authorTemplate *template.Template
	promptTemplate *template.Template
	dataTemplate   *template.Template
	terseTemplate  *template.Template
	titleTemplate  *template.Template
	Series         dalle.Series `json:"series"`
	names          []types.Name
}

func NewApp() *App {
	a := App{
		databases: make(map[string][]string),
	}

	var err error
	if a.Series, err = dalle.GetSeries(); err != nil {
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
	logger.Info("Compiled templates")

	return &a
}

func (a App) String() string {
	bytes, _ := json.MarshalIndent(a, "", "  ")
	return string(bytes)
}

func ReportOn(loc, address, ft, value string) {
	path := filepath.Join("./output/", strings.ToLower(loc))
	file.EstablishFolder(path)
	file.StringToAsciiFile(filepath.Join(path, address+"."+ft), value)
}

func (a *App) toLines(db string) ([]string, error) {
	filename := "./databases/" + db + ".csv"
	lines := file.AsciiFileToLines(filename)
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

func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	} else if a.apiKey = os.Getenv("OPENAI_API_KEY"); a.apiKey == "" {
		log.Fatal("No OPENAI_API_KEY key found")
	}
	if err := a.loadNames(); err != nil {
		logger.Panic(err)
	}
}

func (a *App) DomReady(ctx context.Context) {
	if os.Getenv("TB_CMD_LINE") == "true" {
		return
	}
	if a.session.Load() {
		runtime.WindowSetPosition(a.ctx, a.GetSession().X, a.GetSession().Y)
		runtime.WindowSetSize(a.ctx, a.GetSession().Width, a.GetSession().Height)
	}
	runtime.WindowShow(a.ctx)
}

func (a *App) Shutdown(ctx context.Context) {
	if os.Getenv("TB_CMD_LINE") == "true" {
		return
	}
	a.GetSession().X, a.GetSession().Y = runtime.WindowGetPosition(a.ctx)
	a.GetSession().Width, a.GetSession().Height = runtime.WindowGetSize(a.ctx)
	a.GetSession().Y += 38 // TODO: This is a hack to account for the menu bar - not sure why it's needed
	a.GetSession().Save()
}

func (a *App) GetSession() *config.Session {
	return &a.session
}

func (a *App) GetLastRoute() string {
	return a.GetSession().LastRoute
}

func (a *App) SetLastRoute(route string) {
	a.GetSession().LastRoute = route
	a.GetSession().Save()
}

func (a *App) GetLastTab() string {
	return a.GetSession().LastTab
}

func (a *App) SetLastTab(tab string) {
	a.GetSession().LastTab = tab
	a.GetSession().Save()
}

func isContentPolicyViolation(err error) bool {
	var apiErr struct {
		Code    string `json:"code"`
		Message string `json:"message"`
		Param   string `json:"param"`
		Type    string `json:"type"`
	}
	if jsonErr := json.Unmarshal([]byte(err.Error()), &apiErr); jsonErr == nil {
		return apiErr.Code == "content_policy_violation"
	}
	return false
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
					_, err := a.GetImage(address)
					if err == nil {
						return
					}
					if isContentPolicyViolation(err) {
						msg := fmt.Sprintf("Content policy violation, skipping retry for address: %s Error: %s", address, err)
						logger.Error(msg)
						return
					} else if strings.Contains(err.Error(), "seed length is less than 66") {
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
