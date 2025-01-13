javascript-build:
	bun run build

javascript-watch:
	bun run dev

templ-generate:
	templ generate

templ-watch:
	templ generate --proxy="http://localhost:8080" --watch


# build the project for production
build:
	make templ-generate
	make javascript-build
	@go build -o tmp/main main.go


# live reload
live: 
	make -j3 templ-watch javascript-watch live/server 

live/server:
	APP_ENV=development air \
	--build.cmd="go build -o tmp/main main.go && templ generate --notify-proxy" \
	--build.bin="tmp/main" \
	--build.delay="100" \
	--build.include_ext="go" \
	--build.exclude_dir="node_modules,static/build" \
	--misc.clean_on_exit=true


.PHONY: javascript-build javascript-watch templ-watch templ-generate live build
