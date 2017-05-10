package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
	"time"
)

const usage = `
usage:
	concur <data-dir-path>
`

func processFile(filePath string, q string, ch chan []string) {
	//TODO: open the file, scan each line,
	//do something with the word, and write
	//the results to the channel
	f, err := os.Open(filePath) // open file from file system
	if err != nil {
		log.Fatal(err)
	}
	scanner := bufio.NewScanner(f) // scan one line of file, f, at a time
	n := 0                         // initialize number of words read to 0

	for scanner.Scan() { // returns bool

		// n++ // increment # of words read by 1
		// for i := 0; i < 100; i++ { // hash word w/ SHA256 100 times
		// h := sha256.New()
		// h.Write(scanner.Bytes()) // scanner.Bytes() = bytes read from line in file
		// _ = h.Sum(nil)
		// }

		word := scanner.Text()
		if strings.Contains(word, q) {
			matches := append(matches, word)
		}
	}
	f.Close()     // close scanner
	ch <- matches // write slice of words to channel
}

func processDir(dirPath string, q string) {
	//TODO: iterate over the files in the directory
	//and process each, first in a serial manner,
	//and then in a concurrent manner
	fileinfos, err := ioutil.ReadDir(dirPath)
	if err != nil {
		log.Fatal(err)
	}

	ch := make(chan []string, len(fileinfos))
	for _, fi := range fileinfos { // range over file infos
		go processFile(path.Join(dirPath, fi.Name()), q, ch) // use go keyword to enable concurrent processing
	}
	// nWords := 0 // number of words read

	totalMatches := []string{}
	for i := 0; i < len(fileinfos); i++ {
		matches := <-ch
		totalMatches = append(totalMatches, matches...)
	}
	fmt.Println(strings.Join(totalMatches, ", "))
}

func main() {
	if len(os.Args) < 3 { // check for correct user input
		fmt.Println(usage)
		os.Exit(1)
	}

	dir := os.Args[1]
	q := os.Args[2]

	fmt.Printf("processing directory %s...\n", dir)
	start := time.Now()
	processDir(dir, q)
	fmt.Printf("completed in %v\n", time.Since(start))
}
