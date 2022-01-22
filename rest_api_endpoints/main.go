package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/amulcse/routes"

	"github.com/gorilla/mux"
)

func main() {

	r := mux.NewRouter()
	routes.TweetRouter(r)
	http.Handle("/", r)
	fmt.Printf("Sever Listening: 8000")
	log.Fatal(http.ListenAndServe(":8000", r))

}
