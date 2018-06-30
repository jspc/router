package main

import (
	"flag"
	"log"
	"net/http"
)

var (
	domain = flag.String("domain", ".internal", "Domain on which to accept requests")
	listen = flag.String("listen", ":8080", "Listen address")
)

func main() {
	flag.Parse()

	log.Print("Starting proxy")

	docker, err := NewDocker()
	if err != nil {
		panic(err)
	}

	api := NewAPI(docker, *domain)

	http.Handle("/", api)
	log.Fatal(http.ListenAndServe(*listen, nil))
}
