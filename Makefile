binary=launcher
.DEFAULT_GOAL := help

.PHONY: run
run: debug-build ## Build and run binary without arguments
	./$(binary)

.PHONY: build-linux-amd64
build-linux-amd64: ## build
	env GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o $(binary)-linux-amd64

.PHONY: build-darwin-amd64
build-darwin-amd64: ## build
	env GOOS=darwin GOARCH=amd64 go build -ldflags "-s -w" -o $(binary)-darwin-amd64

.PHONY: build-darwin-arm64
build-darwin-arm64: ## build
	env GOOS=darwin GOARCH=arm64 go build -ldflags "-s -w" -o $(binary)-darwin-arm64

.PHONY: build-windows-amd64
build-windows-amd64: ## build
	env GOOS=windows GOARCH=amd64 go build -ldflags "-s -w" -o $(binary)-windows-amd64

.PHONY: install-linux-amd64
install-linux-amd64: build-linux-amd64 ## install
	mv $(binary)-linux-amd64 ~/bin/$(binary)

.PHONY: install-darwin-amd64
install-darwin-amd64: build-darwin-amd64 ## install
	mv $(binary)-darwin-amd64 ~/bin/$(binary)

.PHONY: install-darwin-arm64
install-darwin-arm64: build-darwin-arm64 ## install
	mv $(binary)-darwin-arm64 ~/bin/$(binary)

.PHONY: install-windows-amd64
install-windows-amd64: build-windows-amd64 ## install
	mv $(binary)-windows-amd64 ~/bin/$(binary)

debug-build: ## Build binary
	go build -o $(binary)

test: ## Run tests
	go test ./...

dependencies:
	go get ./...

cover: ## Run test-coverage and open in browser
	go test -v -covermode=count -coverprofile=coverage.out ./... && go tool cover -html=coverage.out

quick-cover: ## Run simple coverage
	go test -cover ./...

fmt: ## Format source-tree
	gofmt -l -s -w .


help: ## Print all available make-commands
	@grep -E '^[a-zA-Z0-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
