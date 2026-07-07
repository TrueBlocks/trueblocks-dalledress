MSG ?= update

.PHONY: build dev start test lint type-check tidy clean clobber add commit push

build:
	yarn build

dev:
	yarn dev

start: dev

test:
	yarn test

lint:
	yarn lint

type-check:
	yarn type-check

tidy:
	go mod tidy

clean:
	rm -rf bin build/bin frontend/dist

clobber: clean
	@find . \( -path './.git' -o -path './.git/*' \) -prune -o -type d -name node_modules -print -exec rm -rf {} +

add:
	@git add -A

commit:
	@git add -A
	@git commit -m "$(MSG)" || true

push:
	@git add -A
	@git commit -m "$(MSG)" || true
	@git push
