package main

import (
	"bufio"
	"flag"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/pressly/chi"
)

var (
	usernameFlag *string
	passwordFlag *string
)

func main() {
	r := chi.NewRouter()

	r.With(basicAuth).Get("/", showFiles)
	r.With(basicAuth).Get("/*", showFiles)
	r.With(statusCodeHandler).Post("/", handlePost)
	r.With(statusCodeHandler).Post("/*", handlePost)

	port := flag.String("port", "80", "Specifies the port to listen for incoming connections")
	useTls := flag.Bool("tls", false, "Tells gobble to listen for secure connections (ie. https)")
	tlsPort := flag.String("tlsPort", "443", "Specifies the port to listen for incoming secure connections")
	tlsCert := flag.String("tlsCert", "cert.pem", "Specifies the path to the x509 certificate")
	tlsKey := flag.String("tlsKey", "key.pem", "Specifies the path to the private key corresponding to the x509 certificate")
	usernameFlag = flag.String("username", "", "Specify a username to protect againt unauthorized reading of your requests")
	passwordFlag = flag.String("password", "", "Specify a password to protect against unauthorized reading of your requests")

	homeDir := flag.String("dir", "public", "Specifies the root directory which all directories and requests will be stored under")
	flag.Parse()

	err := os.MkdirAll(*homeDir, 0644)
	if err != nil {
		panic("unable to create dir")
	}
	os.Chdir(*homeDir)

	if *useTls == true {
		go func(tlsPort *string, tlsCert *string, tlsKey *string) {
			log.Println("Starting secure server on port " + *tlsPort)
			log.Fatal(http.ListenAndServeTLS(":" + *tlsPort, *tlsCert, *tlsKey, r))
		}(tlsPort, tlsCert, tlsKey)
	}

	log.Println("Starting server on port " + *port)
	log.Fatal(http.ListenAndServe(":" + *port, r))
}

func showFiles(w http.ResponseWriter, r *http.Request) {

	t := template.Must(template.New("index").Parse(`{{define "index"}}
		{{range .Files}}
		<a href="{{$.Path}}/{{.Name}}">{{.Name}}</a><br/>
		{{end}}
		{{end}}`))

	path := chi.URLParam(r, "*")

	if info, err := os.Stat("./" + path); err == nil {
		if info.IsDir() {
			files, err := ioutil.ReadDir("./" + path)
			if err != nil {
				panic("Unable to read directory")
			}

			templData := struct {
				Path  string
				Files []os.FileInfo
			}{
				info.Name(),
				files,
			}

			t.Execute(w, templData)
		} else {
			f, _ := ioutil.ReadFile(path)
			w.Write(f)
		}
	} else {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	}
}

func handlePost(w http.ResponseWriter, r *http.Request) {
	t := time.Now()
	dir := r.URL.Query().Get("dir")
	if dir != "" {
		//Make sure the requested directory is around
		dir = dir + "/" + t.Format("2006-01-02")
		err := os.MkdirAll(dir, 0644)
		if err != nil {
			panic("unable to create dir")
		}
	} else {
		//No directory requested so we give them the default
		dir = t.Format("2006-01-02")
		err := os.MkdirAll(dir, 0644)
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

	fmt.Fprintln(writer, "POST ", r.RequestURI)
	//Write headers to file
	for k, v := range r.Header {
		fmt.Fprintln(writer, k+":", strings.Join(v, ","))
	}

	fmt.Fprintln(writer)
	//Write request body to file
	io.Copy(writer, r.Body)
	w.Write([]byte(fo.Name()[1:]))
}

func statusCodeHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("status_code") != "" {
			status, err := strconv.Atoi(r.URL.Query().Get("status_code"))
			if err == nil {
				w.WriteHeader(status)
			} else {
				log.Println("Invalid number in status_code field")
			}
		}
		next.ServeHTTP(w, r)
	})
}

func basicAuth(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, password, _ := r.BasicAuth()

		//If no auth was set up then we just serve the page
		if *usernameFlag == "" || *passwordFlag == "" {
			h.ServeHTTP(w, r)
			return
		}

		//Auth was configured so we check to make sure the user has the correct credentials
		if username != *usernameFlag || password != *passwordFlag {
			w.Header().Add("WWW-Authenticate", `Basic realm="Restricted"`)
			http.Error(w, "Not Authorized", http.StatusUnauthorized)
			return
		}
		h.ServeHTTP(w, r)
	})
}
