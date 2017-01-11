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
)

func main() {
	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Recover())

	// Route => handler
	e.GET("/", handleGet)

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

func handleGet(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!\n")
}