package main

import (
	"fmt"
	"image"
	"image/color"
	"sort"
	"strconv"
	"strings"

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

// func main_annotate() {
// 	path, err := annotate(os.Args[1], "bottom", 0.2)
// 	if err != nil {
// 		fmt.Println("Error annotating image:", err)
// 		return
// 	}
// 	fmt.Println("Annotated image saved to:", path)
// 	utils.System("open " + path)
// }

// func (a *App) GetImprovedPrompt(ensOrAddr string) (string, string) {
// 	apiKey := os.Getenv("OPENAI_API_KEY")
// 	if apiKey == "" {
// 		log.Fatal("No OPENAI_API_KEY key found")
// 	}
// 	// fmt.Println("API key found", apiKey)

// 	prompt := dd.Prompt // a.GetPrompt(ensOrAddr)
// 	prompt = "Please give me a detailed, but terse, rewriting of this description of an image. Be imaginative. " + prompt + " Mute the colors."
// 	url := "https://api.openai.com/v1/chat/completions"
// 	payload := DalleRequest{
// 		Model: "gpt-3.5-turbo",
// 		Seed:  1337,
// 	}
// 	payload.Messages = append(payload.Messages, Message{Role: "system", Content: prompt})

// 	payloadBytes, err := json.Marshal(payload)
// 	if err != nil {
// 		return fmt.Sprintf("Error: %s", err), ""
// 	}

// 	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadBytes))
// 	if err != nil {
// 		return fmt.Sprintf("Error: %s", err), ""
// 	}
// 	req.Header.Set("Content-Type", "application/json")
// 	req.Header.Set("Authorization", "Bearer "+apiKey)

// 	client := &http.Client{}
// 	resp, err := client.Do(req)
// 	if err != nil {
// 		return fmt.Sprintf("Error: %s", err), ""
// 	}
// 	defer resp.Body.Close()

// 	body, err := io.ReadAll(resp.Body)
// 	if err != nil {
// 		return fmt.Sprintf("Error: %s", err), ""
// 	}

// 	return string(body), prompt
// }
