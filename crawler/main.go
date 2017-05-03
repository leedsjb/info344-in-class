package main

import (
	"fmt"
	"os"
	"time"
)

const usage = `
usage:
	crawler <starting-url>
`

func worker(linkq chan string, resultsq chan []string) {
	for link := range linkq {
		plinks, err := getPageLinks(link)
		if err != nil {
			fmt.Printf("ERROR fetching %s: %v\n", link, err)
			continue // exit this loop iteration
		}

		fmt.Printf("%s (%d links)\n", link, len(plinks.Links))
		time.Sleep(time.Millisecond * 500) // sleep for 1/2 a second
		if len(plinks.Links) > 0 {

			go func(links []string) { // create anonymous go routine function
				resultsq <- links
			}(plinks.Links)

			// resultsq <- plinks.Links // issue w/ deadlock situation
		}
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println(usage)
		os.Exit(1)
	}

	nWorkers := 10                  // number of workers
	linkq := make(chan string, 100) // queue for links, buffered w/ 100 slots
	resultsq := make(chan []string)

	for i := 0; i < nWorkers; i++ {
		go worker(linkq, resultsq)
	}

	linkq <- os.Args[1] // pass starting URL from command line to linkqueue

	// map safe because it is in main.go so only one go routine can write at the same time
	seen := map[string]bool{}     // map of string to bool values to track if a link has been seen
	for links := range resultsq { // results queue is where worker writes links found
		for _, link := range links { // ignore numberic index from slice
			if !seen[link] {
				seen[link] = true
				linkq <- link // pass link from worker to link queue
			}
		}
	}
}
