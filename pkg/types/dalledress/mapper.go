package dalledress

import (
	dalle "github.com/TrueBlocks/trueblocks-dalle/v2"
)

// TODO: THIS CAN BE REMOVED. IT'S NOT USED.
// mapGalleryItemToDress converts a GalleryItem to a minimal DalleDress used in page responses
func mapGalleryItemToDress(gi *DalleDress) *dalle.DalleDress {
	if gi == nil {
		return nil
	}
	return &dalle.DalleDress{
		Original:      gi.Original,
		FileName:      gi.FileName,
		AnnotatedPath: gi.AnnotatedPath,
		ImageURL:      gi.ImageURL,
		Series:        gi.Series,
		Completed:     true,
		CacheHit:      true,
	}
}
