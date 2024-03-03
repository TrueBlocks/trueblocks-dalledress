all:
	go build main.go

all2:
	go build app.go chopper.go dalle.go images.go main.go main_annotate.go main_images.go main_stitch.go main_test.go main_wails.go openai.go seeder.go validate.go

