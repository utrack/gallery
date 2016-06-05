package ifaceHttp

import (
	"github.com/utrack/gallery/storage"
	"net/http"
)

// Upload handles the upload HTTP POST request.
func Upload(u storage.Saver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Not found", http.StatusNotFound)
			return
		}

		err := r.ParseMultipartForm(2048576)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		files := r.MultipartForm.File["upload"]

		for i, _ := range files {
			// Open the file's reader
			rdr, err := files[i].Open()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			// Try to save the file
			err = u.Upload(files[i].Filename, rdr)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
		http.Redirect(w, r, r.Referer(), 302)
	}
}
