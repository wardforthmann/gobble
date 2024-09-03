package main

import (
	"bufio"
	"bytes"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
)

func showFiles(w http.ResponseWriter, r *http.Request) {

	t := template.Must(template.New("index").Parse(`{{define "index"}}
	{{range .Files}}
	<a href="{{$.Path}}/{{.Name}}">{{.Name}}</a><br/>
	{{end}}
	{{end}}`))

	path := chi.URLParam(r, "*")

	//Sanitize the path input and then add the '.' to keep the links relative to the working directory
	//This is necessary to keep badly maliciously formatted paths from escaping the working directory
	path = "." + filepath.Clean("/"+path)

	if info, err := os.Stat(path); err == nil {
		if info.IsDir() {
			files, err := os.ReadDir(path)
			if err != nil {
				panic("Unable to read directory")
			}

			templData := struct {
				Path  string
				Files []fs.DirEntry
			}{
				info.Name(),
				files,
			}

			t.Execute(w, templData)
		} else {
			showFile(r, w, path)
		}
	} else {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	}
}

// Adds the file to the response. A noHeader query parameter will cause the headers to be stripped
// prior to adding them to the response.
func showFile(r *http.Request, w http.ResponseWriter, fileName string) {
	if exists := r.URL.Query().Get("no_header"); exists != "" {
		fileBytes, _ := os.ReadFile(fileName)
		splitBytes := bytes.SplitN(fileBytes, []byte("\n\n"), 2)

		scanner := bufio.NewScanner(bytes.NewReader(splitBytes[0]))

		// Fetch the content-type from the headers if it's available and add it to the response
		contentType := "text/plain"
		for scanner.Scan() {
			header := bytes.Split(scanner.Bytes(), []byte(": "))
			if bytes.Equal(header[0], []byte("Content-Type")) {
				contentType = string(header[1])
				break
			}
		}
		w.Header().Add("Content-Type", contentType)
		w.Write(splitBytes[1])
	} else {
		fileBytes, _ := os.ReadFile(fileName)
		w.Header().Add("Content-Type", "text/plain")
		w.Write(fileBytes)
	}
}

func handlePost(w http.ResponseWriter, r *http.Request) {
	t := time.Now()
	dir := r.URL.Query().Get("dir")
	if dir != "" {
		//Make sure the requested directory is around
		dir = dir + "/" + t.Format("2006-01-02")
		err := os.MkdirAll(dir, 0744)
		if err != nil {
			panic("unable to create dir")
		}
	} else {
		//No directory requested so we give them the default
		dir = t.Format("2006-01-02")
		err := os.MkdirAll(dir, 0744)
		if err != nil {
			panic("unable to create dir")
		}
	}

	//Create file which is named after the create time
	fo, err := os.Create("./" + dir + "/" + t.Format("15.04.05.0000"))
	if err != nil {
		panic(err)
	}
	// close fo on exit and check for its returned error
	defer func() {
		if err := fo.Close(); err != nil {
			panic(err)
		}
	}()

	writer := bufio.NewWriter(fo)
	defer writer.Flush()

	fmt.Fprintln(writer, "POST ", r.URL.Path)
	//Write headers to file
	for k, v := range r.Header {
		fmt.Fprintln(writer, k+":", strings.Join(v, ","))
	}

	fmt.Fprintln(writer)
	//Write request body to file
	io.Copy(writer, r.Body)
	w.Write([]byte(fo.Name()[1:]))
}

func statusCodeHandler() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if status_code := r.URL.Query().Get("status_code"); status_code != "" {
				status, err := strconv.Atoi(status_code)
				if err == nil {
					w.WriteHeader(status)
				} else {
					log.Println("Invalid number in status_code field")
				}
			}
			next.ServeHTTP(w, r)
		})
	}
}
