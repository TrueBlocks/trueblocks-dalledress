.PHONY: app

update:
	@go get "github.com/TrueBlocks/trueblocks-sdk/v5@latest"
	@go get github.com/TrueBlocks/trueblocks-core/src/apps/chifra@latest
	@go mod tidy
	@cd frontend ; yarn upgrade --latest ; cd -

lint:
	@yarn lint

test:
	@export $(grep -v '^#' ../.env | xargs) >/dev/null && yarn test

app:
	@rm -fR build/bin
	@wails build
	@open build/bin/TrueBlocks\ Dalledress.app/

clean:
	@rm -fR node_modules
	@rm -fR frontend/node_modules
	@rm -fR build/bin

