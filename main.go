package main

import "github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/logger"

func main() {
	i := Images
	// i := Annotate // Stitch
	// i := Stitch
	// i := Wails

	switch i {
	case Images:
		main_images()
	case Annotate:
		main_annotate()
	case Stitch:
		main_stitch()
	case Wails:
		main_wails()
	default:
		logger.Panic("Invalid mode")
	}
}

type Mode int

const (
	Unused Mode = iota
	Images
	Annotate
	Stitch
	Wails
)

// package main

// import (
// 	"encoding/json"
// 	"fmt"
// 	"os"
// 	"sync"
// 	"time"

// 	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/colors"
// 	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/file"
// 	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/logger"
// 	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/rpc"
// 	"github.com/ethereum/go-ethereum/common/hexutil"
// 	"github.com/ethereum/go-ethereum/crypto"
// )

// type Database struct {
// 	FileName string
// 	Rows     []string
// }

// type App1 struct {
// 	Databases map[string]Database
// 	conn      *rpc.Connection
// }

// type DalleDress struct {
// 	Orig          string   `json:"orig"`
// 	Seed          string   `json:"seed"`
// 	Parts         []string `json:"parts"`
// 	Prompt        string   `json:"prompt"`
// 	Enhanced      string   `json:"enhanced"`
// 	ImageFilename string   `json:"imageFilename"`
// }

// func (d DalleDress) String() string {
// 	bytes, _ := json.MarshalIndent(d, "", "  ")
// 	return string(bytes)
// }

// func main() {
// 	var wg sync.WaitGroup
// 	file.EstablishFolder("newprompts")
// 	validateChan := make(chan string)
// 	seederChan := make(chan DalleDress)
// 	chopperChan := make(chan DalleDress)
// 	builderChan := make(chan DalleDress)
// 	promptChan := make(chan DalleDress)
// 	enhancedChan := make(chan DalleDress)
// 	imageChan := make(chan DalleDress)

// 	a := App1{}
// 	if a.conn = rpc.NewConnection("mainnet", true, map[string]bool{
// 		"blocks":       true,
// 		"receipts":     true,
// 		"transactions": true,
// 		"traces":       true,
// 		"logs":         true,
// 		"statements":   true,
// 		"state":        true,
// 		"tokens":       true,
// 		"results":      true,
// 	}); a.conn == nil {
// 		logger.Error("Could not find rpc server.")
// 	}

// 	go a.validateInput(validateChan, seederChan)
// 	go a.seeder(seederChan, chopperChan)
// 	go a.chopper(chopperChan, builderChan)
// 	go a.prompter(builderChan, promptChan)
// 	go a.enhancer(promptChan, enhancedChan)
// 	go a.generator(enhancedChan, imageChan, &wg)

// 	// for bn := uint64(0); bn < 19300000; bn += 19300000 / 10000 {
// 	// 	fmt.Printf("chifra blocks --hashes %d --fmt json\n", bn)
// 	// 	// if block, err := a.conn.GetBlockBodyByNumber(bn); err != nil {
// 	// 	// 	logger.Error(err)
// 	// 	// } else {
// 	// 	// 	fmt.Println(block.Hash.Hex())
// 	// 	// 	wg.Add(1)
// 	// 	// 	go func(ln string) {
// 	// 	// 		validateChan <- ln
// 	// 	// 	}(block.Hash.Hex())
// 	// 	// }
// 	// }
// 	// fn := "addresses.txt"
// 	fn := "blocks.txt"
// 	lines := file.AsciiFileToLines(fn)
// 	for _, line := range lines {
// 		wg.Add(1)
// 		go func(ln string) {
// 			validateChan <- ln
// 		}(line)
// 	}

// 	go func() {
// 		wg.Wait()
// 		close(imageChan)
// 	}()

// 	cnt := 1
// 	for image := range imageChan {
// 		a.LogComplete(image.ImageFilename)
// 		if cnt%5 == 0 {
// 			time.Sleep(250 * time.Millisecond) // to let the other channels clear
// 			a.LogSleeper(sleepBetween)
// 			time.Sleep(sleepBetween)
// 		}
// 		cnt++
// 	}
// }

// func (a *App1) validateInput(input chan string, output chan DalleDress) {
// 	for orig := range input {
// 		if isValid(orig) {
// 			a.LogProgress("validateInput", orig)
// 			output <- DalleDress{Orig: orig}
// 		} else {
// 			fmt.Println("Invalid input")
// 			// Properly handle invalid input without halting the entire process
// 		}
// 	}
// 	close(output)
// }

// func isValid(input string) bool {
// 	// Assume implementation of isValid exists
// 	return true
// }

// func (a *App1) seeder(input chan DalleDress, output chan DalleDress) {
// 	for dress := range input {
// 		hash := hexutil.Encode(crypto.Keccak256([]byte(dress.Orig)))
// 		dress.Seed = hash[2:] + dress.Orig[2:]
// 		a.LogProgress("seeder", dress.Seed)
// 		output <- dress
// 	}
// 	close(output)
// }

// func (a *App1) chopper(input chan DalleDress, output chan DalleDress) {
// 	for dress := range input {
// 		dress.Parts = []string{dress.Seed[:4], dress.Seed[4:8], dress.Seed[8:12]}
// 		a.LogProgress("chopper", dress.Parts)
// 		output <- dress
// 	}
// 	close(output)
// }

// func (a *App1) prompter(input chan DalleDress, output chan DalleDress) {
// 	for dress := range input {
// 		dress.Prompt = fmt.Sprintf("A dress with colors %s, pattern %s, and style %s", dress.Parts[0], dress.Parts[1], dress.Parts[2])
// 		a.LogProgress("prompter", dress.Prompt)
// 		output <- dress
// 	}
// 	close(output)
// }

// func (a *App1) enhancer(input chan DalleDress, output chan DalleDress) {
// 	for dress := range input {
// 		dress.Enhanced = dress.Prompt + " in a beautiful outdoor setting"
// 		a.LogProgress("enhancer", dress.Enhanced)
// 		output <- dress
// 	}
// 	close(output)
// }

// var sleepBetween = (time.Millisecond * 250) // 5000)

// func (a *App1) generator(input chan DalleDress, output chan DalleDress, wg *sync.WaitGroup) {
// 	// cnt := 1
// 	for dress := range input {
// 		filename := fmt.Sprintf("newprompts/%s.txt", dress.Orig)
// 		file, err := os.Create(filename)
// 		if err != nil {
// 			fmt.Println("Error creating file:", err)
// 			dress.ImageFilename = "Error creating file"
// 		} else {
// 			file.WriteString(dress.Enhanced)
// 			file.Close()
// 			dress.ImageFilename = filename
// 		}
// 		output <- dress
// 		wg.Done()
// 		// if cnt%5 == 0 {
// 		// 	a.LogSleeper(sleepBetween)
// 		// 	time.Sleep(sleepBetween)
// 		// }
// 		// cnt++
// 	}
// }

// var m sync.Mutex

// func (a *App1) LogProgress(stage string, message interface{}) {
// 	m.Lock()
// 	defer m.Unlock()
// 	logger.Info(colors.Green+"Progress:", stage, message, colors.Off)
// }

// var nImages = 0

// func (a *App1) LogComplete(message string) {
// 	m.Lock()
// 	defer m.Unlock()
// 	nImages++
// 	logger.Info(fmt.Sprintf(colors.Cyan+"Completed[%d]: %s", nImages, message+colors.Off))
// }

// func (a *App1) LogSleeper(duration time.Duration) {
// 	m.Lock()
// 	defer m.Unlock()
// 	logger.Info(colors.Yellow+"Sleeping for", duration, "seconds", colors.Off)
// }
