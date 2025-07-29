package app

import (
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/preferences"
)

// GetAppPreferences returns the current AppPreferences
func (a *App) GetAppPreferences() *preferences.AppPreferences {
	a.prefsMu.RLock()
	defer a.prefsMu.RUnlock()
	appPrefsCopy := a.Preferences.App
	return &appPrefsCopy
}

// SetAppPreferences sets the AppPreferences
func (a *App) SetAppPreferences(appPrefs *preferences.AppPreferences) error {
	a.prefsMu.Lock()
	defer a.prefsMu.Unlock()

	a.Preferences.App = *appPrefs
	return preferences.SetAppPreferences(appPrefs)
}

// SetHelpCollapsed sets the help collapsed state
func (a *App) SetHelpCollapsed(collapse bool) {
	a.prefsMu.Lock()
	defer a.prefsMu.Unlock()

	a.Preferences.App.HelpCollapsed = collapse
	_ = preferences.SetAppPreferences(&a.Preferences.App)
}

// GetLanguage returns the current language setting
func (a *App) GetLanguage() string {
	a.prefsMu.RLock()
	defer a.prefsMu.RUnlock()
	return a.Preferences.App.LastLanguage
}

// SetLanguage sets the language setting
func (a *App) SetLanguage(language string) {
	a.prefsMu.Lock()
	defer a.prefsMu.Unlock()

	a.Preferences.App.LastLanguage = language
	_ = preferences.SetAppPreferences(&a.Preferences.App)
}

func (a *App) GetActiveProjectPath() string {
	a.prefsMu.RLock()
	defer a.prefsMu.RUnlock()
	if len(a.Preferences.App.RecentProjects) == 0 {
		return ""
	}
	return a.Preferences.App.RecentProjects[0]
}

func (a *App) SetActiveProjectPath(path string) {
	a.prefsMu.Lock()
	defer a.prefsMu.Unlock()
	_ = a.Preferences.AddRecentProject(path)
}

// SetMenuCollapsed sets the menu collapsed state
func (a *App) SetMenuCollapsed(collapse bool) {
	a.prefsMu.Lock()
	defer a.prefsMu.Unlock()

	a.Preferences.App.MenuCollapsed = collapse
	_ = preferences.SetAppPreferences(&a.Preferences.App)
}

// GetTheme returns the current theme setting
func (a *App) GetTheme() string {
	a.prefsMu.RLock()
	defer a.prefsMu.RUnlock()
	return a.Preferences.App.LastTheme
}

// SetTheme sets the theme setting
func (a *App) SetTheme(theme string) {
	a.prefsMu.Lock()
	defer a.prefsMu.Unlock()

	a.Preferences.App.LastTheme = theme
	_ = preferences.SetAppPreferences(&a.Preferences.App)
}

// GetDebugMode returns the current debug mode setting
func (a *App) GetDebugMode() bool {
	a.prefsMu.RLock()
	defer a.prefsMu.RUnlock()
	return a.Preferences.App.DebugMode
}

// SetDebugMode sets the debug mode setting
func (a *App) SetDebugMode(debugMode bool) {
	a.prefsMu.Lock()
	defer a.prefsMu.Unlock()

	a.Preferences.App.DebugMode = debugMode
	_ = preferences.SetAppPreferences(&a.Preferences.App)
}
