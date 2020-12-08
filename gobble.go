package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {

	port := flag.String("port", "80", "Specifies the port to listen for incoming connections")
	useTls := flag.Bool("tls", false, "Tells gobble to listen for secure connections (ie. https)")
	tlsPort := flag.String("tlsPort", "443", "Specifies the port to listen for incoming secure connections")
	tlsCert := flag.String("tlsCert", "cert.pem", "Specifies the path to the x509 certificate")
	tlsKey := flag.String("tlsKey", "key.pem", "Specifies the path to the private key corresponding to the x509 certificate")
	usernameFlag := flag.String("username", "", "Specify a username to protect against unauthorized reading of your requests")
	passwordFlag := flag.String("password", "", "Specify a password to protect against unauthorized reading of your requests")

	homeDir := flag.String("dir", "public", "Specify the root directory which all directories and requests will be stored under")
	flag.Parse()

	r := gin.New()
	r.Use(gin.Recovery())

	if *usernameFlag != "" {
		r.GET("/*path", gin.BasicAuth(gin.Accounts{*usernameFlag: *passwordFlag}), showFiles)
	} else {
		r.GET("/*path", showFiles)

	}
	r.POST("/*path", statusCodeHandler, handlePost)

	err := os.MkdirAll(*homeDir, 0744)
	if err != nil {
		panic("Unable to create home directory: " + err.Error())
	}
	err = os.Chdir(*homeDir)
	if err != nil {
		panic("Unable to switch to home directory: " + err.Error())
	}

	if *useTls == true {
		go func(tlsPort *string, tlsCert *string, tlsKey *string) {
			log.Println("Starting secure server on port " + *tlsPort)
			log.Fatal(r.RunTLS(":"+*tlsPort, *tlsCert, *tlsKey))
		}(tlsPort, tlsCert, tlsKey)
	}

	log.Println("Starting server on port " + *port)
	log.Fatal(r.Run(":" + *port))
}

func showFiles(c *gin.Context) {

	t := template.Must(template.New("index").Parse(`{{define "index"}}
		{{range .Files}}
		<a href="{{$.Path}}/{{.Name}}">{{.Name}}</a><br/>
		{{end}}
		{{end}}`))

	path := c.Param("path")
	//Sanitize the path input and then add the '.' to keep the links relative to the working directory
	//This is necessary to keep badly maliciously formatted paths from escaping the working directory
	path = "." + filepath.Clean("/"+path)

	if info, err := os.Stat(path); err == nil {
		if info.IsDir() {
			files, err := ioutil.ReadDir(path)
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

			t.Execute(c.Writer, templData)
		} else {
			showFile(c, path)
		}
	} else {
		c.String(http.StatusNotFound, http.StatusText(http.StatusNotFound))
	}
}

// Adds the file to the response. A noHeader query parameter will cause the headers to be stripped
// prior to adding them to the response.
func showFile(c *gin.Context, fileName string) {
	if _, exists := c.GetQuery("no_header"); exists {
		fileBytes, _ := ioutil.ReadFile(fileName)
		splitBytes := bytes.SplitN(fileBytes, []byte("\n\n"), 2)

		scanner := bufio.NewScanner(bytes.NewReader(splitBytes[0]))

		// Fetch the content-type from the headers if it's available and add it to the response
		contentType := "text/plain"
		for scanner.Scan() {
			header := bytes.Split(scanner.Bytes(), []byte(": "))
			if bytes.Compare(header[0], []byte("Content-Type")) == 0 {
				contentType = string(header[1])
				break
			}
		}

		c.Data(http.StatusOK, contentType, splitBytes[1])
	} else {
		fileBytes, _ := ioutil.ReadFile(fileName)
		c.Data(http.StatusOK, "text/plain", fileBytes)
	}
}

func handlePost(c *gin.Context) {
	t := time.Now()
	dir := c.Query("dir")
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

	fmt.Fprintln(writer, "POST ", c.FullPath())
	//Write headers to file
	for k, v := range c.Request.Header {
		fmt.Fprintln(writer, k+":", strings.Join(v, ","))
	}

	fmt.Fprintln(writer)
	//Write request body to file
	io.Copy(writer, c.Request.Body)
	c.String(http.StatusOK, fo.Name()[1:])
}

func statusCodeHandler(c *gin.Context) {
	if c.Query("status_code") != "" {
		status, err := strconv.Atoi(c.Query("status_code"))
		if err == nil {
			c.Status(status)
		} else {
			log.Println("Invalid number in status_code field")
		}
	}
	c.Next()
}
