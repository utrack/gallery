package main

import (
	"flag"
	"github.com/utrack/gallery/hub"
	"github.com/utrack/gallery/interface/http"
	"github.com/utrack/gallery/storage"
	"log"
	"net/http"
	"path/filepath"
)

var flagPath = flag.String("path", ".", "path to the gallery's directory")
var flagHttpAddr = flag.String("addr", ":8080", "HTTP server's port")

func main() {
	flag.Parse()

	path, err := filepath.Abs(*flagPath)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Wiring up the storage. Path: %v\n", path)

	lister, err := storage.NewLister(path)
	if err != nil {
		log.Fatal(err)
	}

	notifier, err := storage.NewNotifier(path)
	if err != nil {
		log.Fatal(err)
	}

	uploader, err := storage.NewSaver(path)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Starting up the hub")
	h := hub.NewHub(lister, notifier)

	log.Printf("Starting HTTP server on port %v\n", *flagHttpAddr)

	staticTemplate, err := getTemplate()
	if err != nil {
		log.Fatalf("Error when parsing template: %v", err)
	}
	http.HandleFunc(`/`, serveStatic(staticTemplate))
	http.HandleFunc(`/ws`, ifaceHttp.ServeWs(h))
	http.HandleFunc(`/up`, ifaceHttp.Upload(uploader))
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(path))))

	err = http.ListenAndServe(*flagHttpAddr, nil)
	if err != nil {
		log.Fatalf("Error when starting HTTP server: %v", err)
	}
}
