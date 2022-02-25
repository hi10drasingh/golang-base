package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/julienschmidt/httprouter"
)

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
    fmt.Fprint(w, "Welcome!\n")
}

func Hello(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
    fmt.Fprintf(w, "hello, %s!\n", ps.ByName("name"))
}

func main() {

	// If the file doesn't exist, create it or append to the file
    file, err := os.OpenFile("storage/logs/logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
    if err != nil {
        log.Fatal(err)
    }

    log.SetOutput(file)
	
	router := httprouter.New()
    router.GET("/", Index)
    router.GET("/hello/:name", Hello)

    log.Println("Server started at port 8888 test")
    log.Fatal(http.ListenAndServe(":8888", router))
}