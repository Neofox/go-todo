.PHONY: tailwind-watch
tailwind-watch:
	tailwindcss -i ./static/css/input.css -o ./static/css/style.css --watch

.PHONY: tailwind-build
tailwind-build:
	tailwindcss -i ./static/css/input.css -o ./static/css/style.css

.PHONY: templ-watch
templ-watch:
	templ generate --proxy="http://localhost:8080" --watch -v

.PHONY: templ-generate
templ-generate:
	templ generate

.PHONY: build	
build:
	make templ-generate
	make tailwind-build
	@go build -o tmp/main main.go


# live reload
.PHONY: live
live: 
	make -j4 live/templ live/server live/tailwind live/sync_assets

live/templ:
	make templ-watch

live/server:
	ENV=development air \
	--build.cmd="go build -o tmp/main main.go" \
	--build.bin="tmp/main" \
	--build.delay="100" \
	--build.include_ext="go" \
	--build.exclude_dir="" \
	--misc.clean_on_exit=true

live/sync_assets:
	air \
	--build.cmd="templ generate --notify-proxy" \
	--build.bin="true" \
	--build.delay="100" \
	--build.include_ext="js,css" \
	--build.exclude_dir="" \
	--misc.clean_on_exit=true

live/tailwind:
	make tailwind-watch
