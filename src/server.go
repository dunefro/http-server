package main

import (
	"crypto/sha256"
	"crypto/subtle"
	"fmt"
	"github.com/dunefro/http-server/data"
	"log"
	"net/http"
	"os"
	"time"
)

type application struct {
	auth struct {
		username string
		password string
	}
}

func main() {
	app := application{}
	app.auth.username, app.auth.password = getUsernameAndPassword()

	// a multiplexer which is used to serve requests
	mux := http.NewServeMux()
	mux.HandleFunc("/unprotected", app.unprotectedHandler)
	mux.HandleFunc("/protected", app.basicAuth(app.protectedHandler))
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Up and running !!!")
	})
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "healthz")
	})
	// Creating a server type
	srv := &http.Server{
		Addr:         ":9090",
		Handler:      mux,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
	log.Println("Starting server on", srv.Addr)

	_, ok := os.LookupEnv("ENABLE_SSL")
	// if tls is required
	if ok {
		// create key-pair here
		// err := srv.ListenAndServeTLS("./cert.pem", "./mykey.pem")
		// log.Fatal(err)
		log.Println("Do nothng right now")
	} else {
		err := srv.ListenAndServe()
		log.Fatal(err)

	}
}

func (a *application) unprotectedHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, fmt.Sprint(data.GetData()))
}

func (a *application) protectedHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "This is protected handler")
}

// handleFunc is a function
// handlerFunc is a http type which is used to make any oridnary function handle HTTP Request
// the function singature should be func(ResponseWriter, *Request)

func (app *application) basicAuth(hf http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()
		if ok {
			usernameHash := sha256.Sum256([]byte(username))
			passwordHash := sha256.Sum256([]byte(password))
			expectedUsernameHash := sha256.Sum256([]byte(app.auth.username))
			expectedPasswordHash := sha256.Sum256([]byte(app.auth.password))

			userMatch := (subtle.ConstantTimeCompare(usernameHash[:], expectedUsernameHash[:]) == 1)
			passwordMatch := (subtle.ConstantTimeCompare(passwordHash[:], expectedPasswordHash[:]) == 1)

			if userMatch && passwordMatch {
				hf.ServeHTTP(w, r)
				return
			}
		}
		w.Header().Set("WWW-Authenticate", `Basic Realm="Unrestricted", charset="UTF-8"`)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	})
}

func getUsernameAndPassword() (username string, password string) {

	username, ok := os.LookupEnv("AUTH_USERNAME")
	if !ok {
		panic("emtpy username. set the value of env variable AUTH_USERNAME")
	}
	password, ok = os.LookupEnv("AUTH_PASSWORD")
	if !ok {
		panic("emtpy password. set the value of env variable AUTH_PASSWORD")
	}
	return
}
