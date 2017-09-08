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
	STATIC_DIR			= "./public"	// folder name
)

var addr = flag.String("addr", ":9000", "http service address")
var netmonSess = sessions.NewCookieStore([]byte("9c3803d77fb840311dfd9dabd01da5e1"))

//var iest []Franchisee

func logRoute(s *url.URL, m string) {
	log.Printf("%s - %s\n", s, m)
}

func notFoundPage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "tmpl/404.gtpl")
}

func badMethodPage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "tmpl/405.gtpl")
}

func unauthorizedPage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "tmpl/401.gtpl")
}

func dashboardHandler(w http.ResponseWriter, r *http.Request) {
	logRoute(r.URL, r.Method)
	if r.URL.Path != "/dashboard" {
		notFoundPage(w, r)
		return
	}

	if r.Method != "GET" {
		badMethodPage(w, r)
		return
	}

	tmpl := template.Must(template.ParseFiles("tmpl/dashboard.gtpl"))
	tmpl.Execute(w, nil)
}

func outletHandler(w http.ResponseWriter, r *http.Request) {
	logRoute(r.URL, r.Method)
	if r.URL.Path != "/outlet" {
		notFoundPage(w, r)
		return
	}

	if r.Method != "GET" {
		badMethodPage(w, r)
		return
	}

	tmpl := template.Must(template.ParseFiles("tmpl/outlet.gtpl"))
	tmpl.Execute(w, nil)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	logRoute(r.URL, r.Method)
	if r.URL.Path != "/" {
		notFoundPage(w, r)
		return
	}

	if r.Method != "GET" {
		badMethodPage(w, r)
		return
	}

	tmpl := template.Must(template.ParseFiles("tmpl/login.gtpl"))
	tmpl.Execute(w, nil)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	logRoute(r.URL, r.Method)
	if r.URL.Path != "/login" {
		notFoundPage(w, r)
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

	if u == "ABLE" && p == "00ABLE00" {
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

func testHandler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("tmpl/test.gtpl")
	t.Execute(w, nil)
}

func test2Handler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("tmpl/test2.gtpl")
	t.Execute(w, nil)
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
			unauthorizedPage(w, r)
			return
		}
		h.ServeHTTP(w, r)
	})
}

func main() {

	flag.Parse()
	hub := newHub()
	// create room based on operator
	hub.createRoom()
	// run our main server
	go hub.run()

	// setup our routes
	mux := mux.NewRouter()
	// route to the static folder (css/js/img)
	mux.PathPrefix(STATIC_PATH).Handler(http.StripPrefix(STATIC_PATH, http.FileServer(http.Dir(STATIC_DIR))))
	// application routes
	mux.HandleFunc("/", homeHandler)
	mux.HandleFunc("/login", loginHandler)
	mux.HandleFunc("/logout", logoutHandler)
	mux.HandleFunc("/dashboard", dashboardHandler)
	mux.HandleFunc("/outlet", outletHandler)
	mux.HandleFunc("/test", testHandler)
	mux.HandleFunc("/test2", test2Handler)

	// websocket route
	mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})

	if err := http.ListenAndServe(*addr, mux); err != nil {
		log.Fatal("Error starting: ", err)
	}
}