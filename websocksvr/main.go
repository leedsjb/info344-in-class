package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

//HandlerContext provides the HTTP handlers with
//shared values and interfaces
type HandlerContext struct {
	Notifier *Notifier
}

//MessageEvent represents an event with a message
//and a timestamp
type MessageEvent struct {
	Message   string    `json:"message"`
	CreatedAt time.Time `json:"createdAt"`
}

//TriggerEvent triggers a new MessageEvent. This is just a handy
//way to create new events for demo purposes. In a real app, you
//would create and broacast events in response to various handler
//actions, e.g., new user sign-up, post of a new message, etc.
func (ctx *HandlerContext) TriggerEvent(w http.ResponseWriter, r *http.Request) {

	fmt.Println("TriggerEventhandler Called")

	//CORS headers to allow cross-origin requests
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Request-Method", "POST")
	w.Header().Add("Access-Control-Request-Headers", "Content-Type")

	//TODO: create a new MessageEvent with a hard-coded message
	//and the current time for CreatedAt
	//Then pass the MessageEvent to the `.Notify()` method of your notifier
	//so that the event gets broadcasted to all web socket clients

	mE := &MessageEvent{
		Message:   "web sockets are cool",
		CreatedAt: time.Now(),
	}

	// trigger broadcast of new message event to all web socket clients
	ctx.Notifier.Notify(mE)

}

var ugprader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true }, // override Gorilla cross origin check by always returning true
}

//WebSocketUpgradeHandler handles websocket upgrade requests
func (ctx *HandlerContext) WebSocketUpgradeHandler(w http.ResponseWriter, r *http.Request) {
	//TODO: upgrade this request to a web socket connection
	//see https://godoc.org/github.com/gorilla/websocket#hdr-Overview
	//NOTE that by default, the websocket package will reject
	//cross-origin upgrade requests, so make sure you set the
	//CheckOrigin field of the Upgrader to allow upgrades from
	//any origin.
	//See https://godoc.org/github.com/gorilla/websocket#hdr-Origin_Considerations

	conn, err := ugprader.Upgrade(w, r, nil) // TODO *** repsonse header???
	if err != nil {
		http.Error(w, "error upgrading connection"+err.Error(), http.StatusInternalServerError)
	}

	//after upgrading, use the `.AddClient()` method on your notifier
	//to add the new client to your notifier's map of clients

	ctx.Notifier.AddClient(conn)

}

func main() {
	addr := "localhost:4000"

	ctx := &HandlerContext{
		Notifier: NewNotifier(),
	}

	//TODO: start the notifier by calling
	//its .Start() method on a new goroutine
	go ctx.Notifier.Start()

	http.HandleFunc("/v1/ws", ctx.WebSocketUpgradeHandler)
	http.HandleFunc("/v1/trigger", ctx.TriggerEvent)

	fmt.Printf("listening at %s...\n", addr)
	fmt.Printf("test the server by opening the websockclient/index.html page\n")
	fmt.Printf("in a few different browser tabs\n")
	log.Fatal(http.ListenAndServe(addr, nil))
}
