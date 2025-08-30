package dalledress

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
	"strings"

	dallev2 "github.com/TrueBlocks/trueblocks-dalle/v2"
)

// loadCurrentDressFromSidecars attempts to reconstruct a dalle.DalleDress from existing sidecar files
// without triggering any generation. Returns nil if no relevant sidecars exist.
func loadCurrentDressFromSidecars(series, address string) *dallev2.DalleDress {
	if address == "" || series == "" {
		return nil
	}
	addrLower := strings.ToLower(address)
	root := dallev2.OutputDir()
	filename := sanitizeFilename(addrLower)
	selectorPath := filepath.Join(root, series, "selector", filename+".json")

	var dd *dallev2.DalleDress
	if b, err := os.ReadFile(selectorPath); err == nil {
		tmp := dallev2.DalleDress{}
		if json.Unmarshal(b, &tmp) == nil {
			dd = &tmp
		}
	}

	// If no selector JSON, attempt to assemble prompts directly
	if dd == nil {
		// Minimal reconstruction; seed/attributes intentionally omitted to avoid re-execution logic
		dd = &dallev2.DalleDress{
			Original: address,
			FileName: filename,
		}
		// Attempt to read individual prompt files
		dd.Prompt = readIfExists(filepath.Join(root, series, "prompt", filename+".txt"))
		dd.DataPrompt = readIfExists(filepath.Join(root, series, "data", filename+".txt"))
		dd.TitlePrompt = readIfExists(filepath.Join(root, series, "title", filename+".txt"))
		dd.TersePrompt = readIfExists(filepath.Join(root, series, "terse", filename+".txt"))
		dd.EnhancedPrompt = readIfExists(filepath.Join(root, series, "enhanced", filename+".txt"))
	}

	annotatedPath := filepath.Join(root, series, "annotated", addrLower+".png")
	if _, err := os.Stat(annotatedPath); err == nil {
		dd.AnnotatedPath = annotatedPath
		dd.Completed = true
		dd.CacheHit = true
	}
	return dd
}

// sanitizeFilename mirrors the dalle.validFilename (unexported) logic for our limited needs.
func sanitizeFilename(in string) string {
	replace := []string{"/", "\\", ":", "*", "?", "\"", "<", ">", "|"}
	for _, c := range replace {
		in = strings.ReplaceAll(in, c, "_")
	}
	in = strings.TrimSpace(in)
	for strings.Contains(in, "__") { // collapse doubles produced by replacements
		in = strings.ReplaceAll(in, "__", "_")
	}
	return in
}

func readIfExists(p string) string {
	b, err := os.ReadFile(p)
	if err != nil {
		return ""
	}
	return string(b)
}

// selectLatestGalleryItem returns the latest modified gallery item for address (across series) and a sorted slice for thumbnails.
func selectLatestGalleryItem(all []*DalleDress, address string) (latest *DalleDress, sorted []*DalleDress) {
	addr := strings.ToLower(address)
	filtered := make([]*DalleDress, 0, 32)
	for _, gi := range all {
		if gi.Original == addr {
			filtered = append(filtered, gi)
		}
	}
	sort.Slice(filtered, func(i, j int) bool {
		if filtered[i].ModifiedAt == filtered[j].ModifiedAt {
			return filtered[i].FileName < filtered[j].FileName
		}
		return filtered[i].ModifiedAt > filtered[j].ModifiedAt // newest first
	})
	if len(filtered) > 0 {
		latest = filtered[0]
	}
	return latest, filtered
}
