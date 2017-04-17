package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-redis/redis"

	"strings"

	"golang.org/x/net/html"
)

const defaultPort = "80"
const headerContentType = "Content-Type"
const contentTypeHTML = "text/html"
const contentTypeJSON = "application/json; charset=utf-8"

//PageSummary contains summary information about a web page
type PageSummary struct {
	Title string   `json:"title"` //page title
	Links []string `json:"links"` //slice of page URLs
}

//getPageSummary fetches PageSummary info for a given URL
func getPageSummary(URL string) (*PageSummary, error) {
	resp, err := http.Get(URL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error resonse status code: %d", resp.StatusCode)
	}
	if !strings.HasPrefix(resp.Header.Get(headerContentType), contentTypeHTML) {
		return nil, fmt.Errorf("the URL did not return an HTML page")
	}

	psum := &PageSummary{}
	tokenizer := html.NewTokenizer(resp.Body)
	for {
		ttype := tokenizer.Next()
		if ttype == html.ErrorToken {
			return psum, tokenizer.Err()
		}

		//if this is a start tag token
		if ttype == html.StartTagToken {
			token := tokenizer.Token()
			//if this is the page title
			if token.Data == "title" {
				tokenizer.Next()
				psum.Title = tokenizer.Token().Data
			}

			//if this is a hyperlink
			if token.Data == "a" {
				//get the href attribute
				for _, attr := range token.Attr {
					//ignore bookmark links
					if attr.Key == "href" && !strings.HasPrefix(attr.Val, "#") {
						psum.Links = append(psum.Links, attr.Val)
					}
				} //for all attributes
			} //if <a>
		} //if start tag
	} //for each token
} //getPageSummary()

type HandlerContext struct { // context struct contains all fields needed
	redisClient *redis.Client
}

//SummaryHandler handles the /v1/summary resource
func (ctx *HandlerContext) SummaryHandler(w http.ResponseWriter, r *http.Request) { // (ctx *HandlerContext) is a receiver type, gives access to vars in function
	URL := r.FormValue("url")
	if len(URL) == 0 {
		http.Error(w, "please supply a `url` query string parameter", http.StatusBadRequest)
		return
	}

	jbuf, err := ctx.redisClient.Get(URL).Bytes()
	if err != nil && err != redis.Nil { // redis.Nil appears when you look up a key and it's not found
		http.Error(w, "error getting from cache: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if err == redis.Nil {
		//TODO: call getPageSummary() passing URL
		//marshal struct into JSON, and write it
		//to the response

		pgsum, err := getPageSummary(URL)
		if err != nil && err != io.EOF {
			http.Error(w, "Error retrieving page summary: "+err.Error(), http.StatusInternalServerError)
			return
		}

		jbuf, err := json.Marshal(pgsum)
		if err != nil {
			http.Error(w, "Error marshaling JSON: "+err.Error(), http.StatusInternalServerError)
			return
		}

		ctx.redisClient.Set(URL, jbuf, time.Second*60) // save for one minute
	}

	w.Header().Add(headerContentType, contentTypeJSON)
	w.Write(jbuf)
}

func main() {
	host := os.Getenv("HOST")
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = defaultPort
	}
	addr := host + ":" + port

	ropts := redis.Options{
		Addr: "localhost:6379",
	}

	rclient := redis.NewClient(&ropts) // ampersand retrieves address of the following variable
	hctx := &HandlerContext{           // create struct and retrieve address on heap
		redisClient: rclient,
	}

	http.HandleFunc("/v1/summary", hctx.SummaryHandler) // makes HandlerContext available to SummaryHandler

	fmt.Printf("listening at %s...\n", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
