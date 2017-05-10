package main

import (
	"fmt"
	"sync"

	"github.com/gorilla/websocket"
)

//Notifier represents a web sockets notifier
type Notifier struct {
	eventq  chan interface{}
	clients map[*websocket.Conn]bool // boolean value is meaningless, we use a map for fast lookup access of clients
	mu      sync.RWMutex
	//TODO: add other fields you might need
	//such as another channel or a mutex
	//(either would work)
	//remember that go maps ARE NOT safe for
	//concurrent access, so you must do something
	//to protect the `clients` map
}

//NewNotifier constructs a new Notifer.
func NewNotifier() *Notifier {
	//TODO: create, initialize and return
	//a Notifier struct

	bufferSize := 10

	notifier := &Notifier{
		eventq:  make(chan interface{}, bufferSize), // make channel buffered with second param
		clients: make(map[*websocket.Conn]bool),     // make new map w/ websocket.Conn as keys and bools as values (bools are meaningless)
		mu:      sync.RWMutex{},
	}

	return notifier
}

//Start begins a loop that checks for new events
//and broadcasts them to all web socket clients.
//This function should be called on a new goroutine
//e.g., `go mynotifer.Start()`
func (n *Notifier) Start() {
	//TODO: implement this function w/ never ending for loop
	//this should check for new events written
	//to the `eventq` channel, and broadcast
	//them to all of the web socket clients

	for {
		fmt.Println("broadcast loop running")
		event := <-n.eventq
		n.broadcast(event)
	}
}

//AddClient adds a new web socket client to the Notifer
func (n *Notifier) AddClient(client *websocket.Conn) {
	//TODO: implement this
	//But remember that this will be called from
	//an HTTP handler, and each HTTP request is
	//processed on its own goroutine, so your
	//implementation here MUST be safe for concurrent use

	n.mu.Lock() // use mutex write lock to ensure map is safe for concurrent use
	defer n.mu.Unlock()

	clientEntry := n.clients[client]
	if clientEntry == false { // go maps return the "0" value of values that don't exist at a given key
		n.clients[client] = true
	}

}

//Notify will add a new event to the event queue
func (n *Notifier) Notify(event interface{}) {
	//TODO: add the `event` to the `eventq`
	n.eventq <- event // write event to eventq channel
}

//readPump will read all messages (including control messages)
//sent by the client and ignore them. This is necessary in order to
//process the control messages. If you don't do this, the
//websocket will get stuck and start producing errors.
//see https://godoc.org/github.com/gorilla/websocket#hdr-Control_Messages
func (n *Notifier) readPump(client *websocket.Conn) {
	//TODO: implement this according to the notes in the
	//Control Message section of the Gorilla Web Socket docs:
	//https://godoc.org/github.com/gorilla/websocket#hdr-Control_Messages

	for {
		if _, _, err := client.NextReader(); err != nil { // if err is not nil, client has "wandered off" -> close socket
			client.Close()
			fmt.Printf("error reading w/ client.NextReader(): %v\n", err.Error())
			break
		}
	}
}

//broadcast sends the event to all clients as a JSON-encoded object
func (n *Notifier) broadcast(event interface{}) {
	//TODO: Loop over all of the web socket clients in
	//n.clients and write the `event` parameter to the client
	//as a JSON-encoded object.
	//HINT: https://godoc.org/github.com/gorilla/websocket#Conn.WriteJSON
	//and for even better performance, try using a PreparedMessage:
	//https://godoc.org/github.com/gorilla/websocket#PreparedMessage
	//https://godoc.org/github.com/gorilla/websocket#Conn.WritePreparedMessage

	n.mu.Lock()
	defer n.mu.Unlock()

	for client := range n.clients { // iterate over each client in map
		if err := client.WriteJSON(event); err != nil { // write event to client
			client.Close()            // close connection to client
			delete(n.clients, client) // delete client from map TODO ** safe for concurrent access ?? **
		}
	}

	//If you get an error while writing to a client,
	//the client has wandered off, so you should call
	//the `.Close()` method on the client, and delete
	//it from the n.clients map
}
