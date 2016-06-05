package main

import (
	"github.com/utrack/gallery/assets"
	"html/template"
)

func getTemplate() (*template.Template, error) {
	buf, err := assets.Asset("assets/static/gall.tmpl")
	if err != nil {
		return nil, err
	}
	return template.New("gall").Parse(string(buf))

}
