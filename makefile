all:
	go build app.go main.go main_test.go \
		new_app.go new_series.go new_attribute.go new_dalledress.go new_prompts.go \
		new_pipe_addrs.go new_pipe_image.go new_pipe_prompt.go new_pipe_select.go

dev:
	yarn dev
