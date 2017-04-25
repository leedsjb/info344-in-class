package main

import (
	"net/http"
)

//Adapter is the type for adapter functions.
//An adapter function accepts an http.Handler
//and returns a new http.Handler that wraps the
//input handler, providing some pre- and/or
//post-processing.
type Adapter func(http.Handler) http.Handler

//TODO: write an Adapt() function that accepts:
// - handler http.Handler the handler to adapt
// - a variadic slice of Adapter functions
//iterate the slice of Adapter functions in
//reverse order, passing the `handler` to
//each, and resetting `handler` to the
//handler returned from the Adapter func

// Adapt Adapter function
func Adapt(handler http.Handler, adapters ...Adapter) /*can name return here*/ http.Handler { // variadic argument, one or many Adapters can be passed in
	for i := len(adapters) - 1; i >= 0; i-- {
		handler = adapters[i](handler) // passing the handler to each??
	}
	// could also return here regularly, but ABT names the return type so the return is not needed
	return handler
}
