package config

import (
	"encoding/json"
	"path/filepath"

	"github.com/TrueBlocks/trueblocks-browse/pkg/paths"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/file"
)

// Session stores ephemeral things such as last window position, last view, and recent file
type Session struct {
	X      int `json:"x"`
	Y      int `json:"y"`
	Width  int `json:"width"`
	Height int `json:"height"`
}

var defaultSession = Session{
	Width:  1024,
	Height: 768,
}

func NewSession() Session {
	return defaultSession
}

// Load loads the session from the configuration folder. If the file contains
// data, we return true. False otherwise.
func (s *Session) Load() bool {
	fn := getSessionFn()
	if contents := file.AsciiFileToString(fn); len(contents) > 0 {
		if err := json.Unmarshal([]byte(contents), s); err != nil {
			*s = defaultSession
		}
		return true
	} else {
		return false
	}
}

// Save saves the session to the configuration folder.
func (s *Session) Save() {
	fn := getSessionFn()
	if contents, _ := json.Marshal(s); len(contents) > 0 {
		file.StringToAsciiFile(fn, string(contents))
	}
}

// getSessionFn returns the session file name.
func getSessionFn() string {
	if configDir, err := paths.GetConfigDir(); err != nil {
		return "./session.json"
	} else {
		return filepath.Join(configDir, "session.json")
	}
}
