package main

import (
	"embed"
	"fmt"
	"html/template"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

//go:embed templates/*
var content embed.FS

//go:embed static/style.css
var styleCSS []byte

//go:embed static/favicon.ico
var faviconICO []byte

//go:embed static/manifest.json
var manifestJSON []byte

//go:embed static/sw.js
var serviceWorkerJS []byte

//go:embed static/icon-192.png
var icon192PNG []byte

//go:embed static/icon-512.png
var icon512PNG []byte

type Entry struct {
	ID       string
	Content  string
	Type     string
	Filename string
}

func generateUniqueFilename(baseDir, baseName string) string {
	// First try without random prefix
	if _, err := os.Stat(filepath.Join(baseDir, baseName)); os.IsNotExist(err) {
		return baseName
	}
	// If file exists, add random prefix until we find a unique name
	for {
		randChars := fmt.Sprintf("%04d", rand.Intn(10000))
		newName := fmt.Sprintf("%s-%s", randChars, baseName)
		if _, err := os.Stat(filepath.Join(baseDir, newName)); os.IsNotExist(err) {
			return newName
		}
	}
}

// migrateExistingContent moves existing files to the new directory structure
// ONLY REQUIRED FOR MIGRATING FROM THE ORIGINAL VERSION
// TODO: Remove this function in future versions
func migrateExistingContent() error {
	if err := os.MkdirAll(filepath.Join("data", "text"), 0755); err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Join("data", "files"), 0755); err != nil {
		return err
	}
	// Check for existing files
	files, err := os.ReadDir("data")
	if err != nil {
		return err
	}
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		name := file.Name()
		oldPath := filepath.Join("data", name)
		var newPath string
		if strings.HasPrefix(name, "text-") {
			newPath = filepath.Join("data", "text", strings.TrimPrefix(name, "text-"))
		} else if strings.HasPrefix(name, "file-") {
			newPath = filepath.Join("data", "files", strings.TrimPrefix(name, "file-"))
		} else {
			// If there isn't a match, don't do anything
			// This shouldn't happen though, given the implementation
			continue
		}
		if err := os.Rename(oldPath, newPath); err != nil {
			return err
		}
	}
	return nil
}

func main() {
	if err := migrateExistingContent(); err != nil {
		log.Printf("Migration error: %v", err)
	}
	if err := os.MkdirAll("data", 0755); err != nil {
		log.Fatal(err)
	}
	tmpl := template.Must(template.ParseFS(content, "templates/*.html"))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		entries := []Entry{}
		textFiles, _ := os.ReadDir(filepath.Join("data", "text"))
		for _, file := range textFiles {
			if file.IsDir() {
				continue
			}
			data, err := os.ReadFile(filepath.Join("data", "text", file.Name()))
			if err != nil {
				continue
			}
			entries = append(entries, Entry{
				ID:       filepath.Join("text", file.Name()),
				Type:     "text",
				Content:  string(data),
				Filename: file.Name(),
			})
		}
		files, _ := os.ReadDir(filepath.Join("data", "files"))
		for _, file := range files {
			if file.IsDir() {
				continue
			}
			entries = append(entries, Entry{
				ID:       filepath.Join("files", file.Name()),
				Type:     "file",
				Filename: file.Name(),
			})
		}
		tmpl.ExecuteTemplate(w, "index.html", entries)
	})

	http.HandleFunc("/style.css", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/css")
		w.Write(styleCSS)
	})

	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/x-icon")
		w.Write(faviconICO)
	})

	http.HandleFunc("/manifest.json", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write(manifestJSON)
	})

	http.HandleFunc("/sw.js", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/javascript")
		w.Write(serviceWorkerJS)
	})

	http.HandleFunc("/icon-192.png", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/png")
		w.Write(icon192PNG)
	})

	http.HandleFunc("/icon-512.png", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/png")
		w.Write(icon512PNG)
	})

	http.HandleFunc("/submit", func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseMultipartForm(2 << 28); err != nil { // 256 MB
			http.Error(w, err.Error(), 500)
			return
		}
		entryType := r.FormValue("type")
		switch entryType {
		case "text":
			content := r.FormValue("content")
			if content == "" {
				http.Redirect(w, r, "/", http.StatusSeeOther)
				return
			}
			filename := r.FormValue("filename")
			if filename == "" {
				filename = time.Now().Format("Jan-02-15-04-05")
			} else {
				filename = generateUniqueFilename("data/text", filename)
			}

			err := os.WriteFile(filepath.Join("data/text", filename), []byte(content), 0644)
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}
		case "file":
			if err := r.ParseMultipartForm(2 << 28); err != nil {
				http.Error(w, err.Error(), 500)
				return
			}
			files := r.MultipartForm.File["multifile"]
			if len(files) == 0 {
				http.Error(w, "No files uploaded", 400)
				return
			}
			for _, fileHeader := range files {
				if err := func() error {
					file, err := fileHeader.Open()
					if err != nil {
						return fmt.Errorf("failed to open uploaded file: %v", err)
					}
					defer file.Close()
					fileName := generateUniqueFilename("data/files", fileHeader.Filename)
					f, err := os.Create(filepath.Join("data/files", fileName))
					if err != nil {
						return fmt.Errorf("failed to create file: %v", err)
					}
					defer f.Close()
					if _, err := io.Copy(f, file); err != nil {
						return fmt.Errorf("failed to save file: %v", err)
					}
					return nil
				}(); err != nil {
					http.Error(w, err.Error(), 500)
					return
				}
			}
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
	})

	http.HandleFunc("/rename/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		oldPath := strings.TrimPrefix(r.URL.Path, "/rename/")
		newName := r.FormValue("newname")
		if newName == "" {
			http.Error(w, "New name cannot be empty", http.StatusBadRequest)
			return
		}
		baseDir := filepath.Dir(filepath.Join("data", oldPath))
		newName = generateUniqueFilename(baseDir, newName)
		err := os.Rename(
			filepath.Join("data", oldPath),
			filepath.Join(baseDir, newName),
		)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
	})

	http.HandleFunc("/show/", func(w http.ResponseWriter, r *http.Request) {
		id := strings.TrimPrefix(r.URL.Path, "/show/")
		content, err := os.ReadFile(filepath.Join("data", id))
		if err != nil {
			http.Error(w, "File not found", 404)
			return
		}
		w.Header().Set("Content-Type", "text/plain")
		w.Write(content)
	})

	// http.HandleFunc("/download/", func(w http.ResponseWriter, r *http.Request) {
	// 	filename := strings.TrimPrefix(r.URL.Path, "/download/")
	// 	filePath := filepath.Join("data", filename)
	// 	// Set headers to force download
	// 	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%q", filepath.Base(filename)))
	// 	w.Header().Set("Content-Type", "application/octet-stream")
	// 	http.ServeFile(w, r, filePath)
	// })

	http.HandleFunc("/download/", func(w http.ResponseWriter, r *http.Request) {
		filename := strings.TrimPrefix(r.URL.Path, "/download/")
		filePath := filepath.Join("data", filename)
		file, err := os.Open(filePath)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer file.Close()
		// 512 bytes to detect mime-type
		buffer := make([]byte, 512)
		_, err = file.Read(buffer)
		if err != nil && err != io.EOF {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		contentType := http.DetectContentType(buffer)
		// Reset file pointer
		file.Seek(0, 0)
		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%q", filepath.Base(filename)))
		w.Header().Set("Content-Type", contentType)
		http.ServeFile(w, r, filePath)
	})

	http.HandleFunc("/view/", func(w http.ResponseWriter, r *http.Request) {
		filename := strings.TrimPrefix(r.URL.Path, "/view/")
		http.ServeFile(w, r, filepath.Join("data", filename))
	})

	http.HandleFunc("/delete/", func(w http.ResponseWriter, r *http.Request) {
		filename := strings.TrimPrefix(r.URL.Path, "/delete/")
		os.Remove(filepath.Join("data", filename))
		http.Redirect(w, r, "/", http.StatusSeeOther)
	})

	http.HandleFunc("/edit/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		id := strings.TrimPrefix(r.URL.Path, "/edit/")
		if !strings.HasPrefix(id, "text/") {
			http.Error(w, "Can only edit text snippets", http.StatusBadRequest)
			return
		}
		content := r.FormValue("content")
		if content == "" {
			http.Error(w, "Content cannot be empty", http.StatusBadRequest)
			return
		}
		err := os.WriteFile(filepath.Join("data", id), []byte(content), 0644)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
