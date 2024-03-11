package main

import (
	"fmt"
)

var ErrInvalidAddress = fmt.Errorf("not a valid address")
var ErrInvalidSeed = fmt.Errorf("invalid seed")

func (a *App) Chopper(input string) (string, map[string]string, error) {
	_, seed, err := a.SeedBuilder(input)
	if err != nil {
		return "", map[string]string{}, err
	}

	keys := []string{
		"adverb",
		"adjective",
		"noun",
		"emotion",
		"emotionshort",
		"artstyle",
		"artstyle2",
		"litstyle",
		"color1",
		"color2",
		"color3",
		"background",
		"orientation",
	}
	starts := []int{0, 12, 24, 36, 48, 60, 72, 84, 92, 80, 68, 56, 44, 32, 20, 8}
	ends := []int{12, 24, 36, 48, 60, 72, 84, 96, 104, 92, 80, 68, 56, 44, 32, 20}

	segments := make(map[string]string, len(keys))
	for i, key := range keys {
		segments[key] = seed[starts[i]:ends[i]]
	}
	return seed, segments, nil
}
