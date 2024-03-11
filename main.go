package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/colors"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/file"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/logger"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/rpc"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

type Database struct {
	FileName string
	Rows     []string
}

type App1 struct {
	Databases map[string]Database
	conn      *rpc.Connection
}

type DalleDress struct {
	Num           int      `json:"num"`
	Orig          string   `json:"orig"`
	Seed          string   `json:"seed"`
	Parts         []string `json:"parts"`
	Prompt        string   `json:"prompt"`
	Enhanced      string   `json:"enhanced"`
	ImageFilename string   `json:"imageFilename"`
}

func (d DalleDress) String() string {
	bytes, _ := json.MarshalIndent(d, "", "  ")
	return string(bytes)
}

var sleepBetween = 50 * time.Millisecond

// var sleepBetween = 12 * time.Second

func mainshit() {
	var wg sync.WaitGroup
	file.EstablishFolder("newprompts")
	validateChan := make(chan *DalleDress)
	seederChan := make(chan *DalleDress)
	chopperChan := make(chan *DalleDress)
	builderChan := make(chan *DalleDress)
	promptChan := make(chan *DalleDress)
	enhancedChan := make(chan *DalleDress)
	imageChan := make(chan *DalleDress, 5)

	a := App1{}
	if a.conn = rpc.NewConnection("mainnet", true, map[string]bool{
		"blocks":       true,
		"receipts":     true,
		"transactions": true,
		"traces":       true,
		"logs":         true,
		"statements":   true,
		"state":        true,
		"tokens":       true,
		"results":      true,
	}); a.conn == nil {
		logger.Error("Could not find rpc server.")
	}

	go a.validate(validateChan, seederChan)
	go a.seeder(seederChan, chopperChan)
	go a.chopper(chopperChan, builderChan)
	go a.prompter(builderChan, promptChan)
	go a.enhancer(promptChan, enhancedChan)
	go a.generator(enhancedChan, imageChan, &wg)

	fn := "blocks.txt"
	lines := file.AsciiFileToLines(fn)
	for i, line := range lines {
		wg.Add(1)
		go func(orig string) {
			dd := DalleDress{Num: i, Orig: orig}
			validateChan <- &dd
		}(line)
	}

	go func() {
		wg.Wait()
		close(imageChan)
	}()

	for dd := range imageChan {
		a.LogComplete(dd, "completed")
	}
}

func (a *App1) validate(validateChan chan *DalleDress, seederChan chan *DalleDress) {
	for dd := range validateChan {
		if isValid(dd.Orig) {
			// a.LogProgress(dd, "building the seed")
			seederChan <- dd
		} else {
			fmt.Println("Invalid input")
		}
	}
	close(seederChan)
}

func isValid(input string) bool {
	// Assume implementation of isValid exists
	return true
}

func (a *App1) seeder(seederChan chan *DalleDress, chopperChan chan *DalleDress) {
	for dd := range seederChan {
		hash := hexutil.Encode(crypto.Keccak256([]byte(dd.Orig)))
		dd.Seed = hash[2:] + dd.Orig[2:]
		// a.LogProgress(dd, "chopping the seed")
		chopperChan <- dd
	}
	close(chopperChan)
}

func (a *App1) chopper(chopperChan chan *DalleDress, builderChan chan *DalleDress) {
	for dd := range chopperChan {
		dd.Parts = []string{dd.Seed[:4], dd.Seed[4:8], dd.Seed[8:12]}
		// a.LogProgress(dd, "building the prompt")
		builderChan <- dd
	}
	close(builderChan)
}

func (a *App1) prompter(builderChan chan *DalleDress, promptChan chan *DalleDress) {
	for dd := range builderChan {
		dd.Prompt = fmt.Sprintf("A dress with colors %s, pattern %s, and style %s", dd.Parts[0], dd.Parts[1], dd.Parts[2])
		a.LogProgress(dd, "enhancing the prompt")
		promptChan <- dd
	}
	close(promptChan)
}

func (a *App1) enhancer(promptChan chan *DalleDress, enhanceChan chan *DalleDress) {
	for dd := range promptChan {
		dd.Enhanced = dd.Prompt + " in a beautiful outdoor setting"
		enhanceChan <- dd
	}
	close(enhanceChan)
}

var last = time.Now()

func (a *App1) generator(enhanceChan chan *DalleDress, imageChan chan *DalleDress, wg *sync.WaitGroup) {
	for dd := range enhanceChan {
		now := time.Now()
		// since := time.Since(last)
		last = now
		// logger.Warn("Generating image for", dd.Seed[:42], "at", now.Format("2006-01-02 15:04:05"), since)
		filename := fmt.Sprintf("newprompts/%s.txt", dd.Orig)
		file, err := os.Create(filename)
		if err != nil {
			fmt.Println("Error creating file:", err)
			dd.ImageFilename = "Error creating file"
		} else {
			file.WriteString(dd.Enhanced)
			file.Close()
			dd.ImageFilename = filename
		}
		a.LogProgress(dd, "generating the image")
		a.LogSleeper(dd, sleepBetween)
		time.Sleep(sleepBetween)
		imageChan <- dd
		wg.Done()
	}
}

var m sync.Mutex

var nImages = 1
var cols = []string{
	colors.Red,
	colors.Green,
	colors.Yellow,
	colors.Blue,
	colors.Magenta,
	colors.Cyan,
	colors.White,
}

func (a *App1) LogProgress(dd *DalleDress, message string) {
	m.Lock()
	defer m.Unlock()
	if len(dd.Seed) > 42 {
		color := cols[dd.Num%len(cols)]
		color2 := colors.Yellow
		logger.Info(color+dd.Seed[:42], "==>", color2+message, colors.Off)
	}
}

func (a *App1) LogComplete(dd *DalleDress, message string) {
	m.Lock()
	defer m.Unlock()
	color := cols[dd.Num%len(cols)]
	color2 := colors.BrightWhite
	logger.Info(color+dd.Seed[:42], "==>", color2+message+fmt.Sprintf(" %d images", nImages), colors.Off)
	logger.Info(color2+strings.Repeat("=", 70), colors.Off)
	nImages++
}

func (a *App1) LogSleeper(dd *DalleDress, duration time.Duration) {
	m.Lock()
	defer m.Unlock()
	color := cols[dd.Num%len(cols)]
	color2 := colors.Yellow
	logger.Info(color+dd.Seed[:42], "==>", color2+"sleeping for", duration, "seconds", colors.Off)
}
