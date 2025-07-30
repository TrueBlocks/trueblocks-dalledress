package app

import (
	"runtime"

	"github.com/TrueBlocks/trueblocks-dalledress/pkg/msgs"
	"github.com/TrueBlocks/trueblocks-dalledress/pkg/preferences"
	"github.com/wailsapp/wails/v2/pkg/menu"
	"github.com/wailsapp/wails/v2/pkg/menu/keys"
	wailsRuntime "github.com/wailsapp/wails/v2/pkg/runtime"
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

// GetDebugMode returns the current debug mode setting
func (a *App) GetDebugMode() bool {
	a.prefsMu.RLock()
	defer a.prefsMu.RUnlock()
	return a.Preferences.App.DebugMode
}

// SetDebugMode updates the debug mode preference
func (a *App) SetDebugMode(debugMode bool) {
	a.prefsMu.Lock()
	defer a.prefsMu.Unlock()

	a.Preferences.App.DebugMode = debugMode
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

// buildAppMenu creates and configures the application menu structure
func (a *App) buildAppMenu() *menu.Menu {
	appMenu := menu.NewMenu()

	// System Menu (added before File menu)
	system := appMenu.AddSubmenu("System")
	system.AddText("Preferences...", keys.CmdOrCtrl("5"), func(_ *menu.CallbackData) {
		// This matches the cmd+5 keyboard shortcut
		// a.ShowPage("settings")
	})
	system.AddSeparator()
	// TODO: add applicastion name to this menu item
	system.AddText("Quit", keys.CmdOrCtrl("q"), a.FileQuit)

	// File Menu
	file := appMenu.AddSubmenu("File")
	file.AddText("New", keys.CmdOrCtrl("n"), a.FileNew)
	file.AddText("Open", keys.CmdOrCtrl("o"), a.FileOpen)
	file.AddText("Save", keys.CmdOrCtrl("s"), a.FileSave)
	file.AddText("Save As", keys.CmdOrCtrl("shift+s"), a.FileSaveAs)

	if runtime.GOOS == "darwin" {
		appMenu.Append(menu.EditMenu())
	}

	// Window Menu
	window := appMenu.AddSubmenu("Window")
	window.AddText("Minimize", keys.CmdOrCtrl("m"), nil) // menu.WindowMinimize)
	window.AddText("Zoom", nil, nil)                     // menu.WindowZoom)

	// Help Menu
	help := appMenu.AddSubmenu("Help")
	// TODO: add applicastion name to this menu item
	aboutLink := "https://" + preferences.GetAppId().Domain + "/about"
	help.AddText("About", nil, func(_ *menu.CallbackData) {
		wailsRuntime.BrowserOpenURL(a.ctx, aboutLink)
	})
	help.AddText("Report Issue", nil, func(_ *menu.CallbackData) {
		wailsRuntime.BrowserOpenURL(a.ctx, preferences.GetAppId().Github+"/issues")
	})

	return appMenu
}
