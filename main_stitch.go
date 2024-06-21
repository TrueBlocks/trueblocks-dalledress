package main

import (
	"fmt"
	"image"
	"image/color" // or "image/jpeg" if you're using JPEG images
	"log"
	"os"
	"strings"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/base"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/logger"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/utils"
	"github.com/disintegration/imaging"
)

func main_stitch() {
	parts := strings.Split(os.Args[1], "-")
	imagePaths := make([]string, 0, 3)
	start := base.MustParseInt64(os.Getenv("STITCH_START"))
	cnt := base.MustParseInt64(os.Getenv("STITCH_CNT"))
	if cnt == 0 {
		cnt = 3
	}
	for i := start; i < start+cnt; i++ {
		imagePaths = append(imagePaths, parts[0]+"-"+fmt.Sprintf("%05d", i)+".png")
	}

	var images []image.Image
	targetHeight := 500 // set this to your desired height

	// Load and resize images.
	for _, path := range imagePaths {
		img, err := imaging.Open(path)
		if err != nil {
			logger.Error("failed to open image:", err)
			continue
		}

		// Resize the image to the specified height while preserving the aspect ratio.
		resizedImg := imaging.Resize(img, 0, targetHeight, imaging.Lanczos)
		images = append(images, resizedImg)
	}

	// Calculate the total width of the new image (sum of all widths + borders).
	totalWidth := 0
	border := 10 // change border size as needed
	for _, img := range images {
		totalWidth += img.Bounds().Dx() + border
	}
	totalWidth -= border // remove the last border

	// Create a new image with a solid background.
	newImg := imaging.New(totalWidth, targetHeight, color.NRGBA{R: 255, G: 255, B: 255, A: 255})

	offset := 0
	for _, img := range images {
		// Paste the image into the new image at the current offset.
		newImg = imaging.Paste(newImg, img, image.Pt(offset, 0))
		offset += img.Bounds().Dx() + border
	}

	// Save the resulting image.
	stitchedPath := strings.Replace(parts[0]+".png", "annotated/", "stitched/", -1)
	stitchedPath = strings.Replace(stitchedPath, ".png", fmt.Sprintf("-%05d.png", start), -1)
	logger.Info("Stitched image saved to", stitchedPath)
	err := imaging.Save(newImg, stitchedPath)
	if err != nil {
		log.Fatalf("failed to save image: %v", err)
	}
	utils.System("open " + stitchedPath)
}
