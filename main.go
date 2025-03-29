package main

import (
	"fmt"
	"net/http"
)

type api struct {
	addr string
}

func (s *api) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, World!"))
	fmt.Println("Request recieved", r.URL.Path)
	switch r.Method {
	case
		http.MethodGet:
		switch r.URL.Path {
		case "/":
			w.Write([]byte("Index Page"))
			return
		case "/about":
			w.Write([]byte("About Page"))
			return
		}
	default:
		w.Write([]byte("404"))
	}
}

func (a *api) getUsersHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Request recieved", r.URL.Path)
	w.Write([]byte("Users List..."))
}

func (a *api) createUserHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Request recieved", r.URL.Path)
	w.Write([]byte("Users List..."))
}

func main() {
	api := &api{addr: ":8080"}

	mux := http.NewServeMux()

	srv := &http.Server{
		Addr:    api.addr,
		Handler: mux,
	}

	mux.HandleFunc("GET /users", api.getUsersHandler)
	mux.HandleFunc("POST /users", api.createUserHandler)
	err := srv.ListenAndServe()
	if err != nil {
		fmt.Println("Request recieveda")
		panic(err)
	}
}
