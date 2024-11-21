// main.go
package main

import (
	"embed"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gomarkdown/markdown"
	"github.com/google/uuid"
)

//go:embed templates/* static/*
var content embed.FS

//go:embed static/style.css
var styleCSS []byte

type Entry struct {
	ID       string
	Content  string
	Type     string // "text" or "file"
	Filename string // Only used for files
}

func main() {
	// Create data directory if it doesn't exist
	if err := os.MkdirAll("data", 0755); err != nil {
		log.Fatal(err)
	}

	// Parse templates
	tmpl := template.Must(template.ParseFS(content, "templates/*.html"))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		entries := []Entry{}
		files, err := os.ReadDir("data")
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		for _, file := range files {
			if file.Name() == "temporary-print" {
				continue
			}

			data, err := os.ReadFile(filepath.Join("data", file.Name()))
			if err != nil {
				continue
			}

			entry := Entry{ID: file.Name()}
			if strings.HasPrefix(file.Name(), "file-") {
				entry.Type = "file"
				entry.Filename = strings.TrimPrefix(file.Name(), "file-")
			} else {
				entry.Type = "text"
				entry.Content = string(data)
			}
			entries = append(entries, entry)
		}

		tmpl.ExecuteTemplate(w, "index.html", entries)
	})

	http.HandleFunc("/style.css", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/css")
		w.Write(styleCSS)
	})

	http.HandleFunc("/submit", func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseMultipartForm(32 << 20); err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		entryType := r.FormValue("type")
		id := uuid.New().String()

		switch entryType {
		case "text":
			content := r.FormValue("content")
			if content == "" {
				http.Redirect(w, r, "/", http.StatusSeeOther)
				return
			}
			err := os.WriteFile(filepath.Join("data", id), []byte(content), 0644)
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}

		case "file":
			file, header, err := r.FormFile("file")
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}
			defer file.Close()

			f, err := os.Create(filepath.Join("data", "file-"+header.Filename))
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}
			defer f.Close()

			if _, err := io.Copy(f, file); err != nil {
				http.Error(w, err.Error(), 500)
				return
			}
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
	})

	http.HandleFunc("/render/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			id := strings.TrimPrefix(r.URL.Path, "/render/")
			if id != "" {
				content, err := os.ReadFile(filepath.Join("data", id))
				if err != nil {
					http.Error(w, "File not found", 404)
					return
				}
				html := markdown.ToHTML(content, nil, nil)
				tmpl.ExecuteTemplate(w, "render.html", template.HTML(html))
				return
			}
		}

		if r.Method == http.MethodPost {
			content := r.FormValue("content")
			if content == "" {
				http.Redirect(w, r, "/", http.StatusSeeOther)
				return
			}
			html := markdown.ToHTML([]byte(content), nil, nil)
			tmpl.ExecuteTemplate(w, "render.html", template.HTML(html))
			return
		}

		tmpl.ExecuteTemplate(w, "render.html", nil)
	})

	http.HandleFunc("/download/", func(w http.ResponseWriter, r *http.Request) {
		filename := strings.TrimPrefix(r.URL.Path, "/download/")
		http.ServeFile(w, r, filepath.Join("data", filename))
	})

	http.HandleFunc("/delete/", func(w http.ResponseWriter, r *http.Request) {
		filename := strings.TrimPrefix(r.URL.Path, "/delete/")
		os.Remove(filepath.Join("data", filename))
		http.Redirect(w, r, "/", http.StatusSeeOther)
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
