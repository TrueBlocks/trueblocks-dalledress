package app

func (a *App) GetLastRoute() string {
	return a.GetSession().LastRoute
}

func (a *App) GetLastTab() string {
	return a.GetSession().LastTab
}

func (a *App) GetLastAddress() string {
	return a.GetSession().LastAddress
}

func (a *App) GetLastSeries() string {
	return a.GetSession().LastSeries
}

func (a *App) SetLastRoute(route string) {
	a.GetSession().LastRoute = route
	a.GetSession().Save()
}

func (a *App) SetLastTab(tab string) {
	a.GetSession().LastTab = tab
	a.GetSession().Save()
}

func (a *App) SetLastAddress(addr string) {
	a.GetSession().LastAddress = addr
	a.GetSession().Save()
}

func (a *App) SetLastSeries(series string) {
	a.GetSession().LastSeries = series
	a.GetSession().Save()
}
