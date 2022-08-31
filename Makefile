
build:
	rm -rf dist && CGO_ENABLED=0 go build -o ./dist/go_example cmd/main.go
