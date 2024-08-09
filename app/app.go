package app

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand/v2"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"text/template"
	"time"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/base"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/colors"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/file"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/logger"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/output"
	coreTypes "github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/types"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/config"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/daemons"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/dalle"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/messages"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/types"
	"github.com/joho/godotenv"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// Since we need App.ctx to display a dialog and we can only get it when Startup method
// is executed, we keep track of the first fatal error that has happened before Startup
var startupError error

// Find: NewViews
type App struct {
	ctx        context.Context
	Documents  []types.Document
	CurrentDoc *types.Document

	session    config.Session
	apiKeys    map[string]string
	ensMap     map[string]base.Address
	renderCtxs map[base.Address][]*output.RenderCtx
	historyMap map[base.Address]types.SummaryTransaction
	balanceMap map[base.Address]string

	// Summaries
	abis     types.SummaryAbis
	index    types.SummaryIndex
	manifest types.SummaryManifest
	monitors types.SummaryMonitor
	names    types.SummaryName
	status   types.SummaryStatus

	ScraperController *daemons.DaemonScraper
	FreshenController *daemons.DaemonFreshen
	IpfsController    *daemons.DaemonIpfs

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

// Find: NewViews
func NewApp() *App {
	a := App{
		apiKeys:    make(map[string]string),
		renderCtxs: make(map[base.Address][]*output.RenderCtx),
		ensMap:     make(map[string]base.Address),
		// Initialize maps here
		historyMap: make(map[base.Address]types.SummaryTransaction),
		balanceMap: make(map[base.Address]string),
		Documents:  make([]types.Document, 10),
		databases:  make(map[string][]string),
		dalleCache: make(map[string]*dalle.DalleDress),
	}
	a.names.NamesMap = make(map[base.Address]coreTypes.Name)
	a.monitors.MonitorMap = make(map[base.Address]coreTypes.Monitor)
	a.CurrentDoc = &a.Documents[0]
	a.CurrentDoc.Filename = "Untitled"

	// it's okay if it's not found
	_ = a.session.Load()

	if err := godotenv.Load(); err != nil {
		a.Fatal("Error loading .env file")
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

func (a *App) GetContext() context.Context {
	return a.ctx
}

var freshenLock atomic.Uint32

// Freshen gets called by the daemons to instruct first the backend, then the frontend to update.
// Protect against updating too fast... Note that this routine is called as a goroutine.
func (a *App) Freshen(which ...string) {
	// Skip this update we're actively upgrading
	if !freshenLock.CompareAndSwap(0, 1) {
		// logger.Info(colors.Red, "Skipping update", colors.Off)
		return
	}
	logger.Info(colors.Green, "Freshening...", colors.Off)
	defer freshenLock.CompareAndSwap(1, 0)

	notify :=
		func() {
			// Let the front end know it needs to update
			messages.Send(a.ctx, messages.Daemon, messages.NewDaemonMsg(
				a.FreshenController.Color,
				"Freshening...",
				a.FreshenController.Color,
			))
		}

	// First, we want to update the current route if we're told to
	route := ""
	if len(which) > 0 {
		route = which[0]
	}
	switch route {
	case "/abis":
		a.loadAbis(nil)
		notify()
	case "/manifest":
		a.loadManifest(nil)
		notify()
	case "/monitors":
		a.loadMonitors(nil)
		notify()
	case "/names":
		a.loadNames(nil)
		notify()
	case "/index":
		a.loadIndex(nil)
		notify()
	}

	// Now update everything in the fullness of time
	wg := sync.WaitGroup{}
	wg.Add(5)
	go a.loadAbis(&wg)
	go a.loadManifest(&wg)
	go a.loadMonitors(&wg)
	go a.loadNames(&wg)
	go a.loadIndex(&wg)
	wg.Wait()
	notify()
}

// Find: NewViews
func (a *App) Startup(ctx context.Context) {
	a.ctx = ctx

	a.FreshenController = daemons.NewFreshen(a, "freshen", 2000, a.GetLastDaemon("daemon-freshen"))
	a.ScraperController = daemons.NewScraper(a, "scraper", 7000, a.GetLastDaemon("daemon-scraper"))
	a.IpfsController = daemons.NewIpfs(a, "ipfs", 10000, a.GetLastDaemon("daemon-ipfs"))
	go a.startDaemons()

	if startupError != nil {
		a.Fatal(startupError.Error())
	}

	logger.Info("Starting freshen process...")
	a.Freshen(a.GetSession().LastRoute)

	if err := a.loadStatus(); err != nil {
		logger.Panic(err)
	}

	if err := a.loadConfig(); err != nil {
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

func (a *App) Fatal(message string) {
	if message == "" {
		message = "Fatal error occured. The application cannot continue to run."
	}
	log.Println(message)

	// If a.ctx has not been set yet (i.e. we are before calling Startup), we can't display the
	// dialog. Instead, we keep the error and let Startup call this function again when a.ctx is set.
	if a.ctx == nil {
		// We will only display the first error, since it makes more sense
		if startupError == nil {
			startupError = errors.New(message)
		}
		// Return to allow the application to continue starting up, until we get the context
		return
	}
	_, _ = runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
		Type:    runtime.ErrorDialog,
		Title:   "Fatal Error",
		Message: message,
	})
	os.Exit(1)
}
