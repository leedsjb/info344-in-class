package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	addr := "localhost:4000"

	mux := http.NewServeMux() // create your own mux (instead of using default that already exist)
	muxLogged := http.NewServeMux()
	muxLogged.HandleFunc("/v1/hello1", HelloHandler1) // note that the HandleFunc() is being called on the custom mux we created
	muxLogged.HandleFunc("/v1/hello2", HelloHandler2)

	// muxLogged.http.HandleFunc("/v1/hello3", HelloHandler3)

	mux.HandleFunc("/v1/hello3", HelloHandler3)
	logger := log.New(os.Stdout, "", log.LstdFlags)                                            // create a logger that writes to the os.Stdout
	mux.Handle("/v1/", Adapt(muxLogged, logRequests(logger), throttleRequest(2, time.Minute))) // can add additional middleware functions to the 2nd arg in Adapt
	// logRequests(logger)(muxLogged): logRequests(logger) returns a func which is then immediately called using muxLogged as its param

	fmt.Printf("listening at %s...\n", addr)
	log.Fatal(http.ListenAndServe(addr, mux)) // http.DefaultServeMux vs. custom mux

	// http.ListenAndServe uses the default mux when the 2nd argument is nil, otherwise you can use your own custom created mux

}
