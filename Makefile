tailwind-watch:
	tailwindcss -i ./static/css/input.css -o ./static/css/style.css --watch

tailwind-build:
	tailwindcss -i ./static/css/input.css -o ./static/css/style.css --minify

templ-watch:
	templ generate --proxy="http://localhost:8080" --watch

templ-generate:
	templ generate

# build the project for production
build:
	make templ-generate
	make tailwind-build
	@go build -o tmp/main main.go


# live reload
live: 
	make -j4 templ-watch tailwind-watch live/server live/sync_assets

live/server:
	APP_ENV=development air \
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


.PHONY: tailwind-watch tailwind-build templ-watch templ-generate live build
