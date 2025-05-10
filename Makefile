.PHONY: dev css watch-css build clean

dev:
	@echo "Starting development server..."
	air & make watch-css

watch-css:
	@echo "Watching CSS changes..."
	./tailwindcss-linux-x64 -i  ./static/input.css -o ./static/output.css --minify  --watch

css:
	@echo "Building CSS..."
	./tailwindcss-linux-x64 -i ./static/input.css -o ./static/output.css --minify

build:
	@echo "Building application..."
	templ generate
	go build -o ./tmp/main .

clean:
	@echo "Cleaning up..."
	rm -rf tmp/
	rm -f static/output.css 