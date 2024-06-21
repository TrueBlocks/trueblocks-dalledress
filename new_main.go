package main

import (
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/file"
)

func main_new() {
	app := NewApp2()

	addresses := file.AsciiFileToLines("./addresses.txt")
	go app.pipe0_handleAddrs(addresses)
	go app.pipe1_handleSelect()
	go app.pipe2_handlePrompt()
	app.pipe6_handleImage()
}
