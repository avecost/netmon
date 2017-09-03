package main

import (
	"flag"
	"log"
	"net/http"
)

var addr = flag.String("addr", ":9000", "http service address")
//var iest []Franchisee

func serveHome(w http.ResponseWriter, r *http.Request) {

	log.Println(r.URL)
	if r.URL.Path != "/" {
		http.Error(w, "Not found", 404)
		return
	}

	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}

	http.ServeFile(w, r, "/public/index.html")

}

func main() {

	flag.Parse()
	hub := newHub()

	// run our main server
	go hub.run()

	//http.HandleFunc("/", serveHome)
	http.Handle("/able", http.FileServer(http.Dir("./public/")))
	http.Handle("/", http.FileServer(http.Dir("./public/")))
	//http.Handle("/", http.FileServer(http.Dir("/home/whiskie/netmon/public/")))
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})

	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}