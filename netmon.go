package main

import (
	"flag"
	"log"
	"net/http"
	"html/template"
	"net/url"
	"fmt"
	"github.com/gorilla/sessions"
	"time"
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

func postHomeHandler(w http.ResponseWriter, r *http.Request) {
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

	if u == "test" && p == "test" {
		session, _ := netmonSessions.Get(r, "netmon")
		session.Values["user"] = u
		session.Save(r, w)

		http.Redirect(w, r, "/dashboard", 301)
	} else {
		fmt.Fprintf(w, "%s", "Invalid credentials!")
	}
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	sessionOld, err := netmonSessions.Get(r, "netmon")
	fmt.Println("Session in logout")
	fmt.Println(sessionOld)
	if err = sessionOld.Save(r, w); err != nil {
		fmt.Printf("Error saving session: %v", err)
	}
	http.Redirect(w, r, "/", 302)
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

func NoCache(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-cache, private, max-age=0")
		w.Header().Set("Expires", time.Unix(0, 0).Format(http.TimeFormat))
		w.Header().Set("Pragma", "no-cache")
		w.Header().Set("X-Accel-Expires", "0")

		h.ServeHTTP(w, r)
	})
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

	// run our main server
	go hub.run()

	// setup our routes
	mux := http.NewServeMux()
	mux.Handle("/", http.HandlerFunc(homeHandler))
	mux.Handle("/login", http.HandlerFunc(postHomeHandler))
	mux.Handle("/logout", http.HandlerFunc(logoutHandler))
	mux.Handle("/dashboard", NoCache(SecuredRoute(http.HandlerFunc(dashboardHandler))))
	mux.Handle("/sess", NoCache(SecuredRoute(http.HandlerFunc(sessHandler))))

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