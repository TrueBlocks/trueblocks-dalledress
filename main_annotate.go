package main

import (
	"context"
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

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/utils"
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

func contrastColor(cIn color.Color) (color.Color, float64) {
	c, _ := colorful.MakeColor(cIn)
	_, _, l := c.Hcl()
	var contrast colorful.Color
	white := colorful.Color{R: 1, G: 1, B: 1}
	black := colorful.Color{R: 0, G: 0, B: 0}
	if l < 0.5 {
		contrast = white // c.BlendHcl(white, 0.5)
	} else {
		contrast = black // c.BlendHcl(black, 0.5)
	}
	r, g, b := contrast.RGB255()                   // Convert to RGB 0-255 scale
	return color.RGBA{R: r, G: g, B: b, A: 255}, l // Return as color.RGBA with full opacity
}

func getText() string { // path string) string {
	// addr := strings.Replace(strings.Replace(strings.Replace(path, "generated/", "", -1), ".png", "", -1), "/", "", -1)
	// fmt.Println(path, addr)
	// os.Exit(0)
	// parts := strings.Split(addr, "-")
	// addr = parts[0]
	app := NewApp()
	ctx := context.Background()
	app.startup(ctx)
	app.domReady(ctx)
	return "" // app.GetTerse(addr)
}

// annotate reads an image and adds a text annotation to it.
func annotate(fileName, location string, annoPct float64) (ret string, err error) {
	file, err := os.Open(fileName)
	if err != nil {
		return "", err
	}
	defer file.Close()
	img, _, err := image.Decode(file)
	if err != nil {
		return "", err
	}

	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	text := getText() // fileName)
	estimatedFontSize := 20. * (float64(width) / float64(len(text)*7.))

	textWidth := float64(width) * 0.95
	lines := math.Ceil(float64(len(text)) / (textWidth / estimatedFontSize))
	marginHeight := float64(height) * 0.025
	annoHeight := lines * estimatedFontSize * 1.5
	newHeight := height + int(annoHeight+marginHeight*2)

	newImg := image.NewRGBA(image.Rect(0, 0, width, newHeight))
	draw.Draw(newImg, newImg.Bounds(), img, bounds.Min, draw.Src)

	bgColor, _ := findAverageDominantColor(img)
	col, err := parseHexColor(bgColor)
	if err != nil {
		return "", err
	}

	bgRect := image.Rect(0, height, width, newHeight)
	if location == "top" {
		bgRect = image.Rect(0, 0, width, int(annoHeight+marginHeight*2))
		draw.Draw(newImg, bgRect, &image.Uniform{col}, image.Point{}, draw.Src)
	} else {
		draw.Draw(newImg, bgRect, &image.Uniform{col}, image.Point{}, draw.Src)
	}

	gc := gg.NewContextForImage(newImg)
	if err := gc.LoadFontFace("/System/Library/Fonts/Monaco.ttf", estimatedFontSize); err != nil {
		log.Fatalf("Error loading font: %v", err) // Handle the error appropriately
	}
	borderCol := darkenColor(col)
	gc.SetColor(borderCol)
	gc.SetLineWidth(2)
	if location == "top" {
		gc.DrawLine(0, float64(height)*annoPct, float64(width), float64(height)*annoPct)
	} else {
		gc.DrawLine(0, float64(height), float64(width), float64(height))
	}
	gc.Stroke()

	// Draw the text with adjusted margins.
	textColor, _ := contrastColor(col)
	// fmt.Println(l, textColor, col)

	gc.SetColor(textColor) // use the contrasting color for the text
	gc.DrawStringWrapped(text, float64(width)/2, float64(height)+marginHeight*2, 0.5, 0.35, textWidth, 1.5, gg.AlignLeft)

	// Save the new image.
	outputPath := strings.Replace(fileName, "generated/", "annotated/", -1)
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

func main_annotate() {
	path, err := annotate(os.Args[1], "bottom", 0.2)
	if err != nil {
		fmt.Println("Error annotating image:", err)
		return
	}
	fmt.Println("Annotated image saved to:", path)
	utils.System("open " + path)
}
