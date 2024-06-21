package main

func (app *App2) pipe1_handleSelect() {
	for dd := range app.pipe1Chan {
		for i := 0; i < len(dd.Seed); i = i + 8 {
			index := len(dd.Attribs)
			attr := NewAttribute(app, index, dd.Seed[i:i+6])
			dd.Attribs = append(dd.Attribs, attr)
			dd.AttribMap[attr.Name] = attr
			if i+4+6 < len(dd.Seed) {
				index = len(dd.Attribs)
				attr = NewAttribute(app, index, dd.Seed[i+4:i+4+6])
				dd.Attribs = append(dd.Attribs, attr)
				dd.AttribMap[attr.Name] = attr
			}
		}
		app.ReportOn("PostSelect", dd.Orig, dd.String())
		app.pipe2Chan <- dd
	}
	close(app.pipe2Chan)
}
