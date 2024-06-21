package main

func (app *App2) pipe6_handleImage() {
	for dd := range app.pipe6Chan {
		app.ReportDone(dd.Orig)
	}
}
