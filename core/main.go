package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	// set the directory where static files are stored
	dir := "./static/"

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		log.Fatalf("Directory %s does not exist!", dir)
	}

	// serve files from static directory

	fs := http.FileServer(http.Dir(dir))
	http.Handle("/", fs)

	// start the server
	log.Println("serving static files on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}