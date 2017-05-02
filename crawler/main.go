package main

import (
	"fmt"
	"os"
)

const usage = `
usage:
	crawler <starting-url>
`

func worker(toFetch chan string, results chan *PageLinks) {
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println(usage)
		os.Exit(1)
	}

}
