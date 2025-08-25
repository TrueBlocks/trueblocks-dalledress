package dalledress

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	dallev2 "github.com/TrueBlocks/trueblocks-dalle/v2"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/fileserver"
)

type GalleryItem struct {
	Series     string `json:"series"`
	Address    string `json:"address"`
	URL        string `json:"url"`
	RelPath    string `json:"relPath"`
	FileName   string `json:"fileName"`
	Index      int    `json:"index"`
	ModifiedAt int64  `json:"modifiedAt"`
	FileSize   int64  `json:"fileSize"`
}

var addrPattern = regexp.MustCompile(`0x[0-9a-fA-F]{40}`)
var indexPattern = regexp.MustCompile(`(?i)(?:-|_)(\d+)(?:\.png)$`)

func collectGalleryItemsForSeries(root, series string) ([]*GalleryItem, error) {
	if root == "" {
		root = dallev2.OutputDir()
	}
	annotated := filepath.Join(root, series, "annotated")
	entries, err := os.ReadDir(annotated)
	if err != nil {
		return nil, nil
	}
	baseURL := fileserver.CurrentBaseURL()
	items := make([]*GalleryItem, 0, len(entries))
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		if !strings.HasSuffix(strings.ToLower(name), ".png") {
			continue
		}
		address := ""
		if m := addrPattern.FindString(name); m != "" {
			address = strings.ToLower(m)
		}
		index := -1
		if im := indexPattern.FindStringSubmatch(strings.ToLower(name)); len(im) == 2 {
			if v, perr := parsePositiveInt(im[1]); perr == nil {
				index = v
			}
		}
		info, ierr := e.Info()
		modTime := time.Now().Unix()
		size := int64(0)
		if ierr == nil {
			modTime = info.ModTime().Unix()
			size = info.Size()
		}
		relPath := filepath.Join(series, "annotated", name)
		served := strings.ReplaceAll(relPath, string(filepath.Separator), "/")
		url := "file://" + root + "/" + served
		if baseURL != "" {
			url = baseURL + served
		}
		items = append(items, &GalleryItem{Series: series, Address: address, URL: url, RelPath: relPath, FileName: name, Index: index, ModifiedAt: modTime, FileSize: size})
	}
	return items, nil
}

func parsePositiveInt(s string) (int, error) {
	n := 0
	for _, r := range s {
		if r < '0' || r > '9' {
			return 0, ErrInvalidIndex
		}
		n = n*10 + int(r-'0')
		if n > 1000000000 {
			return 0, ErrInvalidIndex
		}
	}
	return n, nil
}

var ErrInvalidIndex = &invalidIndexError{}

type invalidIndexError struct{}

func (e *invalidIndexError) Error() string { return "invalid index" }
