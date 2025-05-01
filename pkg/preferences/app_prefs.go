package preferences

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/kbinani/screenshot"
)

type Bounds struct {
	X      int `json:"x"`
	Y      int `json:"y"`
	Width  int `json:"width"`
	Height int `json:"height"`
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
		Bounds:           getDefaultBounds(),
		LastTab:          make(map[string]string),
		MenuCollapsed:    false,
		HelpCollapsed:    false,
	}
}

func GetAppPreferences() (AppPreferences, error) {
	path := getAppPrefsPath()

	if _, err := os.Stat(path); os.IsNotExist(err) {
		defaults := AppPreferences{
			Version:          "1.0",
			RecentProjects:   []string{},
			LastView:         "/",
			LastViewNoWizard: "/",
			Bounds:           getDefaultBounds(),
			LastTab:          make(map[string]string),
		}

		if err := SetAppPreferences(&defaults); err != nil {
			return AppPreferences{}, err
		}

		return defaults, nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return AppPreferences{}, err
	}

	var appPrefs AppPreferences
	if err := json.Unmarshal(data, &appPrefs); err != nil {
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

func getDefaultBounds() Bounds {
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

func getAppPrefsPath() string {
	return filepath.Join(getConfigBase(), ToCamel(configBaseApp), "app_prefs.json")
}
