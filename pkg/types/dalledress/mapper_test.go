package dalledress

import "testing"

func TestMapGalleryItemToDressNil(t *testing.T) {
	if got := mapGalleryItemToDress(nil); got != nil {
		t.Fatalf("expected nil for nil input, got %#v", got)
	}
}

func TestMapGalleryItemToDressFields(t *testing.T) {
	gi := &DalleDress{
		Series:        "s1",
		Original:      "0xabc",
		ImageURL:      "u",
		AnnotatedPath: "r/p.png",
		FileName:      "p.png",
	}
	dd := mapGalleryItemToDress(gi)
	if dd == nil {
		t.Fatalf("expected non-nil")
	}
	if dd.Original != gi.Original ||
		dd.AnnotatedPath != gi.AnnotatedPath ||
		dd.ImageURL != gi.ImageURL ||
		dd.Series != gi.Series ||
		dd.FileName != gi.FileName {
		t.Fatalf("field mismatch: %#v vs %#v", dd, gi)
	}
	if !dd.Completed || !dd.CacheHit {
		t.Fatalf("expected Completed & CacheHit true")
	}
}
