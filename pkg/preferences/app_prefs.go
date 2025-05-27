package preferences

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/file"
	"github.com/kbinani/screenshot"
)

type Bounds struct {
	X      int `json:"x"`
	Y      int `json:"y"`
	Width  int `json:"width"`
	Height int `json:"height"`
}

func NewBounds() Bounds {
	bounds := screenshot.GetDisplayBounds(0)
	screenW := bounds.Dx()
	screenH := bounds.Dy()

	ret := Bounds{}
	ret.Width = screenW * 3 / 4
	ret.Height = screenH * 3 / 4
	ret.X = (screenW - ret.Width) / 2
	ret.Y = (screenH - ret.Height) / 2
	return ret
}

func (b *Bounds) IsValid() bool {
	return b.X >= 0 && b.Y >= 0 && b.Width > 100 && b.Height > 100
}

type AppPreferences struct {
	Version          string            `json:"version,omitempty"`
	Name             string            `json:"name,omitempty"`
	Bounds           Bounds            `json:"bounds,omitempty"`
	RecentProjects   []string          `json:"recentProjects,omitempty"`
	LastView         string            `json:"lastView,omitempty"`
	LastTab          map[string]string `json:"lastTab,omitempty"`
	LastViewNoWizard string            `json:"lastViewNoWizard,omitempty"`
	MenuCollapsed    bool              `json:"menuCollapsed,omitempty"`
	HelpCollapsed    bool              `json:"helpCollapsed,omitempty"`
}

func (p *AppPreferences) String() string {
	bytes, _ := json.Marshal(p)
	return string(bytes)
}

// NewAppPreferences creates a new AppPreferences instance with default values
func NewAppPreferences() *AppPreferences {
	return &AppPreferences{
		Version:          "1.0",
		RecentProjects:   []string{},
		LastView:         "/",
		LastViewNoWizard: "/",
		Bounds:           NewBounds(),
		LastTab:          make(map[string]string),
		MenuCollapsed:    false,
		HelpCollapsed:    false,
	}
}

func GetAppPreferences() (AppPreferences, error) {
	path := getAppPrefsPath()

	if !file.FileExists(path) {
		defaults := AppPreferences{
			Version:          "1.0",
			RecentProjects:   []string{},
			LastView:         "/",
			LastViewNoWizard: "/",
			Bounds:           NewBounds(),
			LastTab:          make(map[string]string),
		}
		if err := SetAppPreferences(&defaults); err != nil {
			return AppPreferences{}, err
		}

		return defaults, nil
	}

	var appPrefs AppPreferences
	contents := file.AsciiFileToString(path)
	if err := json.Unmarshal([]byte(contents), &appPrefs); err != nil {
		return AppPreferences{}, err
	}

	// Initialize LastTab if nil
	if appPrefs.LastTab == nil {
		appPrefs.LastTab = make(map[string]string)
	}

	return appPrefs, nil
}

func SetAppPreferences(appPrefs *AppPreferences) error {
	path := getAppPrefsPath()

	data, err := json.MarshalIndent(appPrefs, "", "  ")
	if err != nil {
		return err
	}

	err = os.MkdirAll(filepath.Dir(path), 0755)
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

func getAppPrefsPath() string {
	return filepath.Join(getConfigBase(), ToCamel(configBaseApp), "app_prefs.json")
}
