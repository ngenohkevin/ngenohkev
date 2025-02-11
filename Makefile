
templ:
	@templ generate -watch -proxy=http://localhost:3050

tailwind:
	@tailwindcss -i ./static/css/input.css -o ./static/css/styles.css --watch