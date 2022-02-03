package main

import (
	"basicauth/data"
	"crypto/sha256"
	"crypto/subtle"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type application struct {
	auth struct {
		username string
		password string
	}
}

func main() {
	// fmt.Println(data.GetData())
	app := application{}

	app.auth.username, app.auth.password = getUsernameAndPassword()
	// app.auth.password = os.Getenv("AUTH_PASSWORD")

	if app.auth.username == "" {
		log.Fatal("Don't you have a name or what?")
	}
	if app.auth.password == "" {
		log.Fatal("How can I possibly authenticate you without a password?")
	}

	// a multiplexer which is used to server requests
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

	// For https uncomment below
	// err := srv.ListenAndServeTLS("./cert.pem", "./mykey.pem")\

	// For http
	err := srv.ListenAndServe()
	log.Fatal(err)
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

func getUsernameAndPassword() (string, string) {

	data, err := os.ReadFile("./creds.txt")
	if err != nil {
		panic(err.Error)
	}
	var username, password string
	for _, str := range strings.Split(string(data), "\n") {
		kv := strings.Split(str, "=")
		if kv[0] == "AUTH_USERNAME" {
			username = kv[1]
		} else if kv[0] == "AUTH_PASSWORD" {
			password = kv[1]
		} else {
			fmt.Errorf("key %s is invalid", kv[0])

		}
	}
	return username, password
}
