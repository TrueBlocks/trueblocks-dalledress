package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"log"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/file"
	"github.com/fogleman/gg"
	"github.com/lucasb-eyer/go-colorful"
)

// darkenColor slightly darkens a given color.
func darkenColor(c color.Color) color.Color {
	r, g, b, a := c.RGBA()
	factor := 0.9
	return color.RGBA{
		R: uint8(float64(r) * factor),
		G: uint8(float64(g) * factor),
		B: uint8(float64(b) * factor),
		A: uint8(a),
	}
}

func parseHexColor(s string) (color.Color, error) {
	c, err := strconv.ParseUint(strings.TrimPrefix(s, "#"), 16, 32)
	if err != nil {
		return nil, err
	}

	return color.RGBA{
		R: uint8(c >> 16),
		G: uint8(c >> 8 & 0xFF),
		B: uint8(c & 0xFF),
		A: 0xFF,
	}, nil
}

func findAverageDominantColor(img image.Image) (string, error) {
	colorFrequency := make(map[colorful.Color]int)
	bounds := img.Bounds()

	// Count the frequency of each color.
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			col, ok := colorful.MakeColor(img.At(x, y))
			if !ok {
				return "", fmt.Errorf("failed to parse color at %d, %d", x, y)
			}
			colorFrequency[col]++
		}
	}

	// Find the top three most frequent colors.
	type kv struct {
		Key   colorful.Color
		Value int
	}

	var ss []kv
	for k, v := range colorFrequency {
		ss = append(ss, kv{k, v})
	}

	sort.Slice(ss, func(i, j int) bool {
		return ss[i].Value > ss[j].Value
	})

	// Take the top three colors if available.
	topColors := make([]colorful.Color, 0, 3)
	for i := 0; i < len(ss) && i < 3; i++ {
		topColors = append(topColors, ss[i].Key)
	}

	// Average the top three colors.
	var r, g, b float64
	for _, col := range topColors {
		tr, tg, tb := col.RGB255()
		r += float64(tr)
		g += float64(tg)
		b += float64(tb)
	}

	count := float64(len(topColors))
	avgColor := colorful.Color{R: r / count / 255, G: g / count / 255, B: b / count / 255}

	return avgColor.Hex(), nil
}

func contrastColor(bgColor color.Color) color.Color {
	r, g, b, _ := bgColor.RGBA()
	luminance := 0.2126*float64(r) + 0.7152*float64(g) + 0.0722*float64(b)
	if luminance > 0xffff/2 {
		return color.RGBA{0, 0, 0, 0xff} // dark color for light backgrounds
	}
	return color.RGBA{0xff, 0xff, 0xff, 0xff} // light color for dark backgrounds
}

// annotate reads an image and adds a text annotation to it.
func annotate(fileName, location string, annoPct float64) (ret string, err error) {
	text := file.AsciiFileToString(strings.Replace(strings.Replace(fileName, ".png", ".txt", -1), "generated", "txt-generated", -1))
	// Open the image file.
	file, err := os.Open(fileName)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// Decode the image.
	img, _, err := image.Decode(file)
	if err != nil {
		return "", err
	}

	// Find the average dominant color to use as the background color.
	bgColor, _ := findAverageDominantColor(img)

	// Determine the new image dimensions.
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	// Create a new graphics context.
	gc := gg.NewContextForImage(img)

	// Adjust font size based on the length of the text and image width.
	const baseFontSize float64 = 20 // Base font size, adjust as needed
	var estimatedFontSize float64
	if len(text) > 0 {
		estimatedFontSize = baseFontSize * (float64(width) / float64(len(text)*10)) // Adjust the multiplier as needed
		if estimatedFontSize < 12 {                                                 // Set a minimum font size
			estimatedFontSize = 12
		}
	} else {
		estimatedFontSize = baseFontSize
	}

	// Load the font.
	const fontPath string = "/System/Library/Fonts/Monaco.ttf" // Replace with the path to your font file
	if err := gc.LoadFontFace(fontPath, estimatedFontSize); err != nil {
		log.Fatalf("Error loading font: %v", err) // Handle the error appropriately
	}

	// Estimate the number of lines and calculate the annotation rectangle height.
	textWidth := float64(width) * 0.9 // 5% margin on each side for width
	lines := math.Ceil(float64(len(text)) / (textWidth / estimatedFontSize))
	marginHeight := float64(height) * 0.05        // 5% margin for height
	annoHeight := lines * estimatedFontSize * 1.5 // 1.5 for line spacing
	newHeight := height + int(annoHeight+marginHeight*2)

	// Create a new image with the calculated annotation space.
	newImg := image.NewRGBA(image.Rect(0, 0, width, newHeight))
	draw.Draw(newImg, newImg.Bounds(), img, bounds.Min, draw.Src)

	// Parse the background color.
	col, err := parseHexColor(bgColor)
	if err != nil {
		return "", err
	}

	// Draw the background rectangle.
	bgRect := image.Rect(0, height, width, newHeight)
	if location == "top" {
		bgRect = image.Rect(0, 0, width, int(annoHeight+marginHeight*2))
		draw.Draw(newImg, bgRect, &image.Uniform{col}, image.Point{}, draw.Src)
	} else {
		draw.Draw(newImg, bgRect, &image.Uniform{col}, image.Point{}, draw.Src)
	}

	// Draw a border between the new rectangle and the original image.
	borderCol := darkenColor(col)
	gc = gg.NewContextForImage(newImg)
	gc.SetColor(borderCol)
	gc.SetLineWidth(2)
	if location == "top" {
		gc.DrawLine(0, float64(height)*annoPct, float64(width), float64(height)*annoPct)
	} else {
		gc.DrawLine(0, float64(height), float64(width), float64(height))
	}
	gc.Stroke()

	// Draw the text with adjusted margins.
	textColor := contrastColor(col)
	gc.SetColor(textColor) // use the contrasting color for the text
	gc.DrawStringWrapped(text, float64(width)/2, float64(height)+marginHeight*2, 0.5, 0.5, textWidth, 1.5, gg.AlignLeft)

	// Save the new image.
	outputPath := strings.Replace(fileName, ".png", "-annotated.png", -1)
	out, err := os.Create(outputPath)
	if err != nil {
		return "", err
	}
	defer out.Close()

	err = png.Encode(out, gc.Image())
	if err != nil {
		return "", err
	}

	return outputPath, nil
}

func main() {
	// Example usage
	path, err := annotate(os.Args[1], "bottom", 0.2)
	if err != nil {
		fmt.Println("Error annotating image:", err)
		return
	}
	fmt.Println("Annotated image saved to:", path)
}

/*

package main

import (
	"context"
	"embed"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/file"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/logger"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	app := NewApp()
	// logger.Info("nArgs: ", len(os.Args))
	if len(os.Args) > 0 {
		ctx := context.Background()
		app.startup(ctx)
		app.domReady(ctx)
		lines := file.AsciiFileToLines("addresses.txt")
		if len(lines) > 0 {
			wg := sync.WaitGroup{}
			for i := 0; i < len(lines); i++ {
				// fmt.Println(lines[i])
				wg.Add(1)
				go doOne(&wg, app, lines[i])
				if (i+1)%5 == 0 {
					wg.Wait()
					logger.Info("Sleeping for 60 seconds")
					time.Sleep(time.Second * 60)
				}
			}
			wg.Wait()
			return
		}

		if len(os.Args) == 2 {
			for i := 0; i < 10; i++ {
				wg := sync.WaitGroup{}
				for j := 0; j < 5; j++ {
					logger.Info("Round", i, "run", j)
					wg.Add(1)
					go doOne(&wg, app, fmt.Sprintf("0x%040x", 10010010+(i*10)+j)) // os.Args[1])
					// SeedBump++
				}
				wg.Wait()
				logger.Info("Sleeping for 60 seconds")
				time.Sleep(time.Second * 60)
			}
		} else {
			wg := sync.WaitGroup{}
			for _, arg := range os.Args[1:] {
				wg.Add(1)
				go doOne(&wg, app, arg)
			}
			wg.Wait()
		}
	} else {
		opts := options.App{
			Title:            app.GetSettings().Title,
			Width:            app.GetSettings().Width,
			Height:           app.GetSettings().Height,
			BackgroundColour: nil,
			Bind:             app.GetTypes(),
			StartHidden:      true,
			OnStartup:        app.startup,
			OnDomReady:       app.domReady,
			OnShutdown:       app.shutdown,
			AssetServer: &assetserver.Options{
				Assets: assets,
			},
		}
		if err := wails.Run(&opts); err != nil {
			println("Error:", err.Error())
		}
	}
}

*/
