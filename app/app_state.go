package app

func (a *App) GetLast(which string) string {
	switch which {
	case "route":
		return a.GetSession().LastRoute
	case "tab":
		return a.GetSession().LastTab
	case "address":
		return a.GetSession().LastAddress
	case "series":
		return a.GetSession().LastSeries
	}
	return "Unknown"
}

func (a *App) SetLast(which, value string) {
	reload := false
	switch which {
	case "route":
		a.GetSession().LastRoute = value
	case "tab":
		a.GetSession().LastTab = value
	case "address":
		a.GetSession().LastAddress = value
	case "series":
		a.GetSession().LastSeries = value
		reload = true
	}
	a.GetSession().Save()
	if reload {
		a.ReloadDatabases()
	}
}
