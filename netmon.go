package main

import (
	"flag"
	"log"
	"net/http"
	"github.com/gorilla/mux"
	"html/template"
	"net/url"
)

const (
	CLIENT_TTL 			= 1				// mark offline if more than (minutes)
	TICKER_SERVER_TIME 	= 1				// run every (second)
	TICKER_ONLINE_TIME	= 15			// run every (second)

	TIME_ZONE 			= "Asia/Manila"			// local datetime
	TIME_FORMAT 		= "2006-01-02 15:04:05"	// how datetime is formatted

	STATIC_DIR 			= "/public/"			// css/js folder
)

var addr = flag.String("addr", ":9000", "http service address")
//var iest []Franchisee

func logRoute(s *url.URL, m string) {
	log.Printf("%s - %s\n", s, m)
}

func getDashboardHandler(w http.ResponseWriter, r *http.Request) {
	logRoute(r.URL, r.Method)
	if r.URL.Path != "/dashboard" {
		http.Error(w, "Not found", 404)
		return
	}

	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}

	tmpl := template.Must(template.ParseFiles("tmpl/dashboard.html"))
	tmpl.Execute(w, nil)
}

func getLoginHandler(w http.ResponseWriter, r *http.Request) {
	logRoute(r.URL, r.Method)
	if r.URL.Path != "/" {
		http.Error(w, "Not found", 404)
		return
	}

	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}

	tmpl := template.Must(template.ParseFiles("tmpl/login.html"))
	tmpl.Execute(w, nil)
}

func postLoginHandler(w http.ResponseWriter, r *http.Request) {
	logRoute(r.URL, r.Method)
	if r.URL.Path != "/login" {
		http.Error(w, "Not found", 404)
		return
	}

	if r.Method == "GET" {
		t, _ := template.ParseFiles("tmpl/login.html")
		t.Execute(w, nil)
		return
	}

	r.ParseForm()
	log.Println("username: ", r.Form["username"])
	log.Println("password: ", r.Form["password"])
}

func NewRouter() *mux.Router {
	r := mux.NewRouter().StrictSlash(true)
	r.PathPrefix(STATIC_DIR).Handler(http.StripPrefix(STATIC_DIR, http.FileServer(http.Dir("."+STATIC_DIR))))

	return r
}

func main() {

	flag.Parse()
	hub := newHub()

	// run our main server
	go hub.run()

	r := NewRouter()
	r.HandleFunc("/", getLoginHandler).Methods("GET")
	r.HandleFunc("/login", postLoginHandler).Methods("POST")
	r.HandleFunc("/dashboard", getDashboardHandler).Methods("GET")

	// handler for websocket
	r.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})

	if err := http.ListenAndServe(*addr, r); err != nil {
		log.Fatal("Error starting: ", err)
	}

	//http.HandleFunc("/", serveHome)
	//http.Handle("/able", http.FileServer(http.Dir("./public/")))
	//http.Handle("/", http.FileServer(http.Dir("./public/")))
	//http.Handle("/", http.FileServer(http.Dir("/home/whiskie/netmon/public/")))
	//http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
	//	serveWs(hub, w, r)
	//})
	//
	//err := http.ListenAndServe(*addr, nil)
	//if err != nil {
	//	log.Fatal("ListenAndServe: ", err)
	//}

}