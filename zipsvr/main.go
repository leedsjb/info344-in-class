package main

import "fmt"
import "os"
import "log"
import "net/http"

// note: capitalized function names are exported in the Go package

// invoked when someone does a request against a resource path
// 2 args: response writer and pointer to request
// note var name on the left, type on the right (reverse order)
// * means sending by reference, as opposed to by value (like java int and Integer) ***??***
func helloHandler(w http.ResponseWriter, r *http.Request){
	
	name := r.URL.Query().Get("name") // retrieve name from URL params , also := declares and names variables simultaneously
	
	// headers must be set before you send the response
	w.Header().Add("Content-Type", "text/plain") // allows one to return headers as well

	w.Write([]byte("hello " + name)) // convert from string to slice of bytes
}

func main(){
	addr := os.Getenv("ADDR") 

	if len(addr) == 0 {
		log.Fatal("please set ADDR environment variable")
	}

	http.HandleFunc("/hello", helloHandler) // registers helloHandler w/ hello resource path

	fmt.Printf("Server is listening at %s...\n",addr)

	log.Fatal(http.ListenAndServe(addr, nil))

}