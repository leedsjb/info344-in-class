package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"os"
)

// go `` allows string that preserve formatting such as line breaks
const usage = `
usage:
	hmac sign|verify <key> <value>
`

func main() {
	if len(os.Args) < 4 || (os.Args[1] != "sign" && os.Args[1] != "verify") {
		fmt.Println(usage)
		os.Exit(1)
	}

	cmd := os.Args[1]   // sign or verify
	key := os.Args[2]   // key
	value := os.Args[3] // value

	switch cmd {

	case "sign": // go automatically breaks at the end of a case as opposed to "falling through"

		v := []byte(value) // create a byte slice of the value to be signed (underlying array automatically created)

		h := hmac.New(sha256.New, []byte(key)) // creates byte slice from string, key and creates hash using sha256
		h.Write(v)
		sig := h.Sum(nil)
		// fmt.Println(sig)

		buf := make([]byte, len(v)+len(sig))
		copy(buf, v)            // copy value to buffer
		copy(buf[len(v):], sig) // copy sig to end of buf
		fmt.Println(base64.URLEncoding.EncodeToString(buf))

	case "verify":

		buf, err := base64.URLEncoding.DecodeString(value)
		if err != nil {
			fmt.Printf("error decoding: %v\n", err)
			os.Exit(1)
		}
		v := buf[:len(buf)-sha256.Size]
		sig := buf[len(buf)-sha256.Size:]

		h := hmac.New(sha256.New, []byte(key))
		h.Write(v)
		sig2 := h.Sum(nil)

		if hmac.Equal(sig, sig2) { // use to prevent time attacks, constant time equals method
			fmt.Println("signature is valid")
		} else {
			fmt.Println("Invalid signature.")
		}

	}

}
