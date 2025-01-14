# Help command to list all available commands
help:
	@echo "Available commands:"
	@echo "  make live          - Start development server with live reload"
	@echo "  make build         - Build the project for production"
	@echo "Setup:"
	@echo "  make rename-module NEW_NAME=your-module-name - Rename the Go module"

javascript-watch:
	@bun run dev

templ-watch:
	@templ generate --proxy="http://localhost:8080" --watch --open-browser=false 

# build the project for production
build:
	@echo "Building project for production..."
	@templ generate
	@bun run build
	@go build -o bin/main main.go
	@echo "Project built successfully. You can now run the binary in the bin directory."

# Development commands
live: 
	@echo "Starting development server..."
	@echo "\033[0;34m➜\033[0m Local server: http://localhost:8080"
	@echo "\033[0;34m➜\033[0m Proxy server: http://localhost:7331 <- hot reload will be done on this port"
	@echo "\033[0;34m➜\033[0m Starting all services..."
	@make -j3 templ-watch javascript-watch live/server 

live/server:
	@APP_ENV=development air \
	--build.cmd="go build -o tmp/main main.go && templ generate --notify-proxy" \
	--build.bin="tmp/main" \
	--build.delay="100" \
	--build.include_ext="go" \
	--build.exclude_dir="node_modules,static/build" \
	--misc.clean_on_exit=true

# Setup commands
rename-module:
	@if [ "$(NEW_NAME)" = "" ]; then \
		echo "Error: Please provide a module name using NEW_NAME=your-module-name"; \
		exit 1; \
	fi
	@echo "Renaming project to $(NEW_NAME)..."
	$(eval CURRENT_MODULE := $(shell grep -m1 "module" go.mod | cut -d ' ' -f 2))
	$(eval PROJECT_NAME := $(shell basename $(NEW_NAME)))
	@echo "Current module: $(CURRENT_MODULE)"
	@echo "Current package name: $(shell jq -r .name package.json)"
	# Update Go module
	@go mod edit -module $(NEW_NAME)
	@find . -type f -name '*.go' -exec perl -pi -e 's|$(CURRENT_MODULE)|$(NEW_NAME)|g' {} \;
	# Update package.json
	@jq '.name = "$(PROJECT_NAME)"' package.json > package.json.tmp && mv package.json.tmp package.json
	@go mod tidy
	@echo "✓ Module renamed from $(CURRENT_MODULE) to $(NEW_NAME)"
	@echo "✓ Package name updated to $(PROJECT_NAME)"

.PHONY: help javascript-watch templ-watch live build rename-module

# Set default target to help
.DEFAULT_GOAL := help