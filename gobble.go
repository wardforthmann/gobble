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
)

func main() {
	// Echo instance
	e := echo.New()

	// Middleware
	//e.Use(middleware.Logger())
	//e.Use(middleware.Recover())

	// Route => handler
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!\n")
	})

	e.POST("/", func(c echo.Context) error {

		t := time.Now()
		fo, err := os.Create(t.Format("2006-01-02 15.04.05"))
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
	})

	// Start server
	e.Logger.Fatal(e.Start(":1323"))
}