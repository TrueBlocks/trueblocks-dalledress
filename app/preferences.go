package app

import (
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/msgs"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/preferences"
)

// GetUserPreferences returns the current user preferences
func (a *App) GetUserPreferences() *preferences.UserPreferences {
	return &a.Preferences.User
}

// SetUserPreferences updates and persists user preferences
func (a *App) SetUserPreferences(userPrefs *preferences.UserPreferences) error {
	a.Preferences.User = *userPrefs
	return preferences.SetUserPreferences(userPrefs)
}

// GetOrgPreferences returns the current organization preferences
func (a *App) GetOrgPreferences() *preferences.OrgPreferences {
	return &a.Preferences.Org
}

// SetOrgPreferences updates and persists organization preferences
func (a *App) SetOrgPreferences(orgPrefs *preferences.OrgPreferences) error {
	a.Preferences.Org = *orgPrefs
	return preferences.SetOrgPreferences(orgPrefs)
}

// GetAppPreferences returns a copy of current application preferences with thread safety
func (a *App) GetAppPreferences() *preferences.AppPreferences {
	a.prefsMu.RLock()
	defer a.prefsMu.RUnlock()
	appPrefsCopy := a.Preferences.App
	return &appPrefsCopy
}

// SetAppPreferences updates and persists application preferences with thread safety
func (a *App) SetAppPreferences(appPrefs *preferences.AppPreferences) error {
	a.prefsMu.Lock()
	defer a.prefsMu.Unlock()

	a.Preferences.App = *appPrefs
	return preferences.SetAppPreferences(appPrefs)
}

// GetLanguage returns the currently selected language
func (a *App) GetLanguage() string {
	a.prefsMu.RLock()
	defer a.prefsMu.RUnlock()
	return a.Preferences.App.LastLanguage
}

// SetLanguage updates the application language preference
func (a *App) SetLanguage(language string) {
	a.prefsMu.Lock()
	defer a.prefsMu.Unlock()

	a.Preferences.App.LastLanguage = language
	_ = preferences.SetAppPreferences(&a.Preferences.App)
}

// GetTheme returns the currently selected theme
func (a *App) GetTheme() string {
	a.prefsMu.RLock()
	defer a.prefsMu.RUnlock()
	return a.Preferences.App.LastTheme
}

// SetTheme updates the application theme preference
func (a *App) SetTheme(theme string) {
	a.prefsMu.Lock()
	defer a.prefsMu.Unlock()

	a.Preferences.App.LastTheme = theme
	_ = preferences.SetAppPreferences(&a.Preferences.App)
}

// GetActiveProjectPath returns the path of the most recently used project
func (a *App) GetActiveProjectPath() string {
	a.prefsMu.RLock()
	defer a.prefsMu.RUnlock()
	if len(a.Preferences.App.RecentProjects) == 0 {
		return ""
	}
	return a.Preferences.App.RecentProjects[0]
}

// SetActiveProjectPath adds a project path to the recent projects list
func (a *App) SetActiveProjectPath(path string) {
	a.prefsMu.Lock()
	defer a.prefsMu.Unlock()
	_ = a.Preferences.AddRecentProject(path)
}

// updateRecentProjects updates the recent projects list with the active project path
func (a *App) updateRecentProjects() {
	activeProject := a.GetActiveProject()
	if activeProject == nil || activeProject.GetPath() == "" {
		return
	}

	path := activeProject.GetPath()

	if err := a.Preferences.AddRecentProject(path); err != nil {
		msgs.EmitError("add recent project failed", err)
		return
	}

	msgs.EmitManager("update_recent_projects")
}

// GetDebugCollapsed returns the current debug mode setting
func (a *App) GetDebugCollapsed() bool {
	a.prefsMu.RLock()
	defer a.prefsMu.RUnlock()
	return a.Preferences.App.DebugCollapsed
}

// SetDebugCollapsed updates the debug mode preference
func (a *App) SetDebugCollapsed(collapse bool) {
	a.prefsMu.Lock()
	defer a.prefsMu.Unlock()

	a.Preferences.App.DebugCollapsed = collapse
	_ = preferences.SetAppPreferences(&a.Preferences.App)
}

// SetHelpCollapsed updates the help panel collapsed state
func (a *App) SetHelpCollapsed(collapse bool) {
	a.prefsMu.Lock()
	defer a.prefsMu.Unlock()

	a.Preferences.App.HelpCollapsed = collapse
	_ = preferences.SetAppPreferences(&a.Preferences.App)
}

// SetMenuCollapsed updates the menu panel collapsed state
func (a *App) SetMenuCollapsed(collapse bool) {
	a.prefsMu.Lock()
	defer a.prefsMu.Unlock()

	a.Preferences.App.MenuCollapsed = collapse
	_ = preferences.SetAppPreferences(&a.Preferences.App)
}
