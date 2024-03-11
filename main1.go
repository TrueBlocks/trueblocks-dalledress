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
