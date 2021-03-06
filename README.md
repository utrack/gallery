# gallery [![GoDoc](https://godoc.org/github.com/utrack/gallery?status.svg)](https://godoc.org/github.com/utrack/gallery)
[![GitHub license](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/utrack/gallery/blob/master/LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/utrack/gallery)](https://goreportcard.com/report/github.com/utrack/gallery)

Sample gallery written using Go.

# Requirements
Golang compiler and tools (v1.5 or later) are required. See the [official Getting Started guide](https://golang.org/doc/install) or your distro's docs for detailed instructions.

# Installation
```
go get -u github.com/utrack/gallery/cmd/gogallery
```
If you're using Go < 1.6 - you need to set envvar `GO15VENDOREXPERIMENT` to `1` before go-getting:
```
GO15VENDOREXPERIMENT=1 go get -u github.com/utrack/gallery/cmd/gogallery
```

# Running
Check that your `PATH` envvar has `$GOPATH\bin` and run the command:
```
gogallery
```

Use flag -path to provide the path to the pictures:
```
gogallery -path /home/u/Pictures
```
HTTP server runs on addr `:8080` by default; use the `-addr` flag to change that.

# Testing
```
go test github.com/utrack/gallery/cmd/gogallery/...
```
Tests are written using the [GoConvey](https://github.com/smartystreets/goconvey) framework. If you have `goconvey` tools installed in your `$PATH`, cd to the project's path and run `goconvey` to use its web interface.

## TODO
- [ ] Generate thumbnails and show them instead of fullsize images
