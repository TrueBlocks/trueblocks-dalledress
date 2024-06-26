all:
	go build -o bin/dalle app.go attribute.go backend.go dalledress.go main.go prompts.go series.go

dev:
	yarn dev
