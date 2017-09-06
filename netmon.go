package main

import (
	"flag"
	"log"
	"net/http"
	"html/template"
	"net/url"
	"fmt"
	"github.com/gorilla/mux"
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
var netmonSess = sessions.NewCookieStore([]byte("9c3803d77fb840311dfd9dabd01da5e1"))

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

func loginHandler(w http.ResponseWriter, r *http.Request) {
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
		sess, _ := netmonSess.Get(r, "netmon")
		sess.Values["user"] = u
		err := sess.Save(r, w)
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}

		http.Redirect(w, r, "/dashboard", 301)
	} else {
		fmt.Fprintf(w, "%s", "Invalid credentials!")
	}
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	sess, err := netmonSess.Get(r, "netmon")
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	sess.Options.MaxAge = -1

	err = sess.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	http.Redirect(w, r, "/", 301)
}

func sessHandler(w http.ResponseWriter, r *http.Request) {
	logRoute(r.URL, r.Method)

	session, _ := netmonSess.Get(r, "netmon")
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
	mux := mux.NewRouter()
	// route to the static folder (css/js/img)
	mux.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("./public"))))

	mux.HandleFunc("/", homeHandler)
	mux.HandleFunc("/login", loginHandler)
	mux.HandleFunc("/logout", logoutHandler)
	mux.HandleFunc("/dashboard", dashboardHandler)

	// websocket route
	mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})

	if err := http.ListenAndServe(*addr, mux); err != nil {
		log.Fatal("Error starting: ", err)
	}
}