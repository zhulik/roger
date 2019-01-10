package web

import (
	"fmt"
	"log"
	"net/http"
)

func rootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello!")
}

func Serve(port int) {
	http.HandleFunc("/info/", apiHandler)
	http.HandleFunc("/", rootHandler)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
