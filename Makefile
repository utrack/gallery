.PHONY: all

all: build

build: bindata
	go build github.com/utrack/gallery/cmd/gogallery

bindata: assets/assets.go

assets/assets.go: assets/static
	go-bindata -pkg assets -o assets/assets.go assets/static/...
