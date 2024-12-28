/*
https://gist.github.com/paulmach/7271283

Serve is a very simple static file server in go
Usage:

	-port="5175": port to serve on
	-dir=".":     the directory of static files to host

Navigating to http://localhost:5175 will display the index.html or directory
listing file.
*/
package main

import (
	"flag"
	"log"
	"net/http"
)

func main() {
	port := flag.String("port", "5175", "port to serve on")
	directory := flag.String("dir", ".", "the directory of static file to host")
	flag.Parse()

	http.Handle("/", http.FileServer(http.Dir(*directory)))

	log.Printf("Serving %s on HTTP port: %s\n", *directory, *port)
	log.Fatal(http.ListenAndServe(":"+*port, nil))
}
