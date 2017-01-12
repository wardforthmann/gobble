package main

import (
	"net/http"

	"github.com/labstack/echo"
	"fmt"
	"os"
	"bufio"
	"io"
	"strings"
	"time"
	"github.com/labstack/echo/middleware"
	"html/template"
	"io/ioutil"
)

func main() {
	// Echo instance
	e := echo.New()

	t := &Template{
		templates: template.Must(template.New("index").Parse(`{{define "index"}}
		{{range .Files}}
		<a href="{{$.Path}}/{{.Name}}">{{.Name}}</a><br/>
		{{end}}
		{{end}}`)),
	}
	e.Renderer = t

	// Middleware
	e.Use(middleware.Recover())

	// Route => handler
	e.GET("/*", handleGet)

	e.POST("/", handlePostPayload)

	// Start server
	e.Logger.Fatal(e.Start(":1323"))
}

func handlePostPayload(c echo.Context) error {
	t := time.Now()
	dir := c.QueryParam("dir")
	if dir != "" {
		err := os.MkdirAll(dir, 0644)
		if err != nil {
			panic("unable to create dir")
		}
	} else {
		dir = t.Format("2006-01-02")
		err := os.MkdirAll(dir, 0644)
		if err != nil {
			panic("unable to create dir")
		}
	}

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

	w := bufio.NewWriter(fo)
	defer w.Flush()

	for k, v := range c.Request().Header {
		fmt.Fprintln(w, k + ":", strings.Join(v, ","))
	}

	fmt.Fprintln(w)
	io.Copy(w, c.Request().Body)

	return c.String(http.StatusOK, t.Format("2006-01-02 15:04:05"))
}

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func handleGet(c echo.Context) error {
	path := c.Param("*")

	if info, err := os.Stat("./" + path); err == nil {
		if info.IsDir() {
			files, err := ioutil.ReadDir("./" + path)
			if err != nil {
				panic("Unable to read directory")
			}

			templData := struct {
				Path string
				Files []os.FileInfo
			}{
				info.Name(),
				files,
			}

			return c.Render(http.StatusOK, "index", templData)
		} else {
			return c.File(path)
		}
	}

	return c.NoContent(http.StatusNotFound)
}