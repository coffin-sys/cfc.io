package main

import (
	"flag"
	"fmt"
	"github.com/coffin-sys/cfc.io/cfc"
	"github.com/gorilla/mux"
	"github.com/labstack/gommon/log"
	"net/http"
)

var baseHost string

func main() {
	flag.StringVar(&baseHost, "host", "cfc.io", "Base Host")
	flag.Parse()

	j := cfc.New(baseHost)
	r := mux.NewRouter()
	r.HandleFunc("/_ws/", j.WebsocketHandler)
	r.PathPrefix("/").HandlerFunc(j.HttpHandler)

	fmt.Println("Server is running on Port 4200")
	log.Fatal(http.ListenAndServe(":4200", r))
}
