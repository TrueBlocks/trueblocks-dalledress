package dalledress

import (
	dalle "github.com/TrueBlocks/trueblocks-dalle/v2"
)

func GetCurrentDressFor(series, address string) *dalle.DalleDress {
	if address == "" {
		return &dalle.DalleDress{}
	}
	if series == "" {
		series = "empty"
	}
	// Always attempt generation (cached path is fast if image exists)
	_, _ = dalle.GenerateAnnotatedImage(series, address, false, 0)
	if pr := dalle.GetProgress(series, address); pr != nil && pr.DalleDress != nil {
		dd := *pr.DalleDress
		if dd.AnnotatedPath != "" {
			dd.Completed = true
			if pr.CacheHit {
				dd.CacheHit = true
			}
		}
		return &dd
	}
	return &dalle.DalleDress{}
}
