package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	if err := run(os.Stdout, os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func run(w io.Writer, args []string) error {

	port := flag.String("port", "8080", "Specifies the port to listen for incoming connections")
	useTls := flag.Bool("tls", false, "Tells gobble to listen for secure connections (ie. https)")
	tlsPort := flag.String("tlsPort", "443", "Specifies the port to listen for incoming secure connections")
	tlsCert := flag.String("tlsCert", "cert.pem", "Specifies the path to the x509 certificate")
	tlsKey := flag.String("tlsKey", "key.pem", "Specifies the path to the private key corresponding to the x509 certificate")
	usernameFlag := flag.String("username", "", "Specify a username to protect against unauthorized reading of your requests")
	passwordFlag := flag.String("password", "", "Specify a password to protect against unauthorized reading of your requests")

	homeDir := flag.String("dir", "public", "Specify the root directory which all directories and requests will be stored under")
	flag.Parse()

	r := chi.NewRouter()
	r.Use(middleware.Recoverer)

	creds := make(map[string]string)
	if *usernameFlag != "" {
		creds[*usernameFlag] = *passwordFlag
	}

	addRoutes(r, creds)

	err := os.MkdirAll(*homeDir, 0744)
	if err != nil {
		panic("Unable to create home directory: " + err.Error())
	}
	err = os.Chdir(*homeDir)
	if err != nil {
		panic("Unable to switch to home directory: " + err.Error())
	}

	if *useTls {
		go func(tlsPort *string, tlsCert *string, tlsKey *string) {
			log.Println("Starting secure server on port " + *tlsPort)
			log.Fatal(http.ListenAndServeTLS(":"+*tlsPort, *tlsCert, *tlsKey, r))
		}(tlsPort, tlsCert, tlsKey)
	}

	log.Println("Starting server on port " + *port)
	log.Fatal(http.ListenAndServe(":"+*port, r))
	return nil
}
