all:
	go build main.go \
		new_app.go new_attribute.go new_dalledress.go new_main.go new_prompts.go \
		new_pipe_addrs.go new_pipe_image.go new_pipe_prompt.go new_pipe_select.go \
		main_annotate.go main_images.go main_stitch.go main_test.go main_wails.go \
		app.go dalle.go images.go openai.go
