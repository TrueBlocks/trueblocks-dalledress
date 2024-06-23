all:
	go build app.go main.go main_test.go main_wails.go \
		new_main.go new_main_test.go new_app.go \
		new_attribute.go new_dalledress.go new_image.go new_prompts.go \
		new_openai_image.go new_openai_prompt.go \
		new_pipe_addrs.go new_pipe_image.go new_pipe_prompt.go new_pipe_select.go
