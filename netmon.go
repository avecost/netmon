package main

import (
	"flag"
	"log"
	"net/http"
	"html/template"
	"net/url"
	"fmt"
	"github.com/gorilla/sessions"
)

const (
	CLIENT_TTL 			= 1				// mark offline if more than (minutes)
	TICKER_SERVER_TIME 	= 1				// run every (second)
	TICKER_ONLINE_TIME	= 15			// run every (second)

	TIME_ZONE 			= "Asia/Manila"			// local datetime
	TIME_FORMAT 		= "2006-01-02 15:04:05"	// how datetime is formatted

	STATIC_PATH 		= "/public/"	// URL css/js folder
	STATIC_DIR			= "public"		// folder name
)

var addr = flag.String("addr", ":9000", "http service address")
var netmonSessions = sessions.NewCookieStore([]byte("9c3803d77fb840311dfd9dabd01da5e1"))

//var iest []Franchisee

func logRoute(s *url.URL, m string) {
	log.Printf("%s - %s\n", s, m)
}

func dashboardHandler(w http.ResponseWriter, r *http.Request) {
	logRoute(r.URL, r.Method)
	if r.URL.Path != "/dashboard" {
		http.Error(w, "Not found", 404)
		return
	}

	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}

	tmpl := template.Must(template.ParseFiles("tmpl/dashboard.gtpl"))
	tmpl.Execute(w, nil)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	logRoute(r.URL, r.Method)
	if r.URL.Path != "/" {
		http.Error(w, "Not found", 404)
		return
	}

	if r.Method != "GET" {
		http.Error(w, "Method not allowed", 405)
		return
	}

	tmpl := template.Must(template.ParseFiles("tmpl/login.gtpl"))
	tmpl.Execute(w, nil)
}

func postHomeHandler(w http.ResponseWriter, r *http.Request) {
	logRoute(r.URL, r.Method)
	if r.URL.Path != "/login" {
		http.Error(w, "Not found", 404)
		return
	}

	if r.Method == "GET" {
		t, _ := template.ParseFiles("tmpl/login.gtpl")
		t.Execute(w, nil)
		return
	}

	// TODO: perform validation of Credentials
	r.ParseForm()
	u := r.FormValue("username")
	p := r.FormValue("password")

	if u == "test" && p == "test" {
		session, _ := netmonSessions.Get(r, "netmon")
		session.Values["user"] = u
		session.Save(r, w)

		http.Redirect(w, r, "/dashboard", 301)
	} else {
		fmt.Fprintf(w, "%s", "Invalid credentials!")
	}
}

func sessHandler(w http.ResponseWriter, r *http.Request) {
	logRoute(r.URL, r.Method)

	session, _ := netmonSessions.Get(r, "netmon")
	u := session.Values["user"]
	fmt.Fprintf(w, "%s", u)
}

func userAllowed(r string, u string) (b bool) {
	// TODO: perform lookup functions for route and user
	fmt.Printf("R: %s U: %s\n", r, u)

	return true
}

func SecuredRoute(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if ! userAllowed(r.URL.Path, "test") {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		h.ServeHTTP(w, r)
	})
}

func main() {

	flag.Parse()
	hub := newHub()

	// run our main server
	go hub.run()

	// setup our routes
	mux := http.NewServeMux()
	mux.Handle("/", http.HandlerFunc(homeHandler))
	mux.Handle("/login", http.HandlerFunc(postHomeHandler))
	mux.Handle("/dashboard", SecuredRoute(http.HandlerFunc(dashboardHandler)))
	mux.Handle("/sess", SecuredRoute(http.HandlerFunc(sessHandler)))

	// websocket route
	mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})

	// where the sites assets are located (js/css/img)
	mux.Handle(STATIC_PATH, http.StripPrefix(STATIC_PATH, http.FileServer(http.Dir(STATIC_DIR))))

	if err := http.ListenAndServe(*addr, mux); err != nil {
		log.Fatal("Error starting: ", err)
	}

}