package main

import (
	"context"
	"testing"
)

func TestChopper(t *testing.T) {
	ctx := context.Background()

	app := NewApp()
	app.startup(ctx)

	tests := []struct {
		name      string
		input     string
		seed      string
		segements map[string]string
		err       error
	}{
		{
			name:      "invalid address - too long",
			input:     "0x12345678901234567890123456789012345678902",
			seed:      "",
			segements: map[string]string{},
			err:       ErrInvalidAddress,
		},
		{
			name:      "invalid address - too short",
			input:     "0x123456789012345678901234567890123456789",
			seed:      "",
			segements: map[string]string{},
			err:       ErrInvalidAddress,
		},
		{
			name:      "invalid address - no hex",
			input:     "1234567890123456789012345678901234567890",
			seed:      "",
			segements: map[string]string{},
			err:       ErrInvalidAddress,
		},
		{
			name:      "invalid address - not ENS",
			input:     "trueblocks",
			seed:      "",
			segements: map[string]string{},
			err:       ErrInvalidAddress,
		},
		{
			name:      "valid address",
			input:     "0xf503017d7baf7fbc0fff7492b751025c6a78179b",
			seed:      "572979b29ccd964cb6c456ea1baec10878f46f176efb3a20bcf262a9f2bae1c8f503017d7baf7fbc0fff7492b751025c6a78179b",
			segements: testChopperMap,
			err:       nil,
		},
		{
			name:      "valid ens",
			input:     "trueblocks.eth",
			seed:      "572979b29ccd964cb6c456ea1baec10878f46f176efb3a20bcf262a9f2bae1c8f503017d7baf7fbc0fff7492b751025c6a78179b",
			segements: testChopperMap,
			err:       nil,
		},
	}

	for _, test := range tests {
		seed, segments, err := app.Chopper(test.input)
		if err != test.err {
			t.Errorf("%s: Error: got %v, want %v", test.name, err, test.err)
		}
		if len(segments) != len(test.segements) {
			t.Errorf("%s: Segments length: got %v, want %v", test.name, len(segments), len(test.segements))
		}
		if len(segments) > 0 {
			cnt := 0
			for key, value := range testChopperMap {
				if segments[key] != value {
					t.Errorf("%s: Segment %d key %s value mismatch: got %s, want %s", test.name, cnt, key, segments[key], value)
				}
				cnt++
			}
		}
		if seed != test.seed {
			t.Errorf("%s: Seed: got %s, want %s", test.name, seed, test.seed)
		}
	}
}

var testChopperMap = map[string]string{
	"adverb":       "572979b29ccd",
	"adjective":    "964cb6c456ea",
	"noun":         "e1c8f503017d",
	"emotion":      "6f176efb3a20",
	"emotionshort": "1baec10878f4",
	"artstyle":     "7baf7fbc0fff",
	"artstyle2":    "7492b751025c",
	"litstyle":     "bcf262a9f2ba",
	"color1":       "025c6a78179b",
	"color2":       "0fff7492b751",
	"color3":       "017d7baf7fbc",
	"background":   "56ea1baec108",
	"orientation":  "9ccd964cb6c4",
}
