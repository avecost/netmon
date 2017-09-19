package main

import (
	"fmt"
	"os"
	"flag"
	"log"
	"strings"
	"net/url"
	"net/http"
	"html/template"

	"github.com/gorilla/mux"
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

// server runtime config
var addr = flag.String("addr", ":9000", "http service address")
// server allowed users
var Users = []AppUser{}

//var netmonSess = sessions.NewCookieStore([]byte("9c3803d77fb840311dfd9dabd01da5e1"))
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
	if UserAllowedURL(Users, Username(r), r.URL.Path) {
		tpl, err := template.ParseFiles("tmpl/dashboard.gtpl")
		if err != nil {
			fmt.Println("Error parsing template")
		}
		tpl.Execute(w, nil)
	} else {
		http.Redirect(w, r, "/", 302)
	}
}

func terminalHandler(w http.ResponseWriter, r *http.Request) {
	type Outlet struct {
		Operator string
		Name string
	}

	r.ParseForm()
	operator := strings.ToUpper(r.FormValue("o"))
	outlet := strings.ToUpper(r.FormValue("t"))
	d := Outlet{Operator: operator, Name: outlet}

	if UserAllowedURL(Users, Username(r), r.URL.Path) {
		tpl, err := template.ParseFiles("tmpl/terminal.gtpl")
		if err != nil {
			fmt.Println("Error parsing template")
		}
		tpl.Execute(w, d)
	} else {
		http.Redirect(w, r, "/", 302)
	}
}

func outletHandler(w http.ResponseWriter, r *http.Request) {
	type Operator struct {
		Name string
	}

	r.ParseForm()
	operator := strings.ToUpper(r.FormValue("o"))
	d := Operator{Name: operator}

	if UserAllowedURL(Users, Username(r), r.URL.Path) {
		tpl, err := template.ParseFiles("tmpl/outlet.gtpl")
		if err != nil {
			fmt.Println("Error parsing template")
		}
		tpl.Execute(w, d)
	} else {
		http.Redirect(w, r, "/", 302)
	}
}

func operatorHandler(w http.ResponseWriter, r *http.Request) {
	type Operator struct {
		Name string
	}

	r.ParseForm()
	operator := strings.ToUpper(r.FormValue("o"))
	d := Operator{Name: operator}

	if UserAllowedURL(Users, Username(r), r.URL.Path) {
		tpl, err := template.ParseFiles("tmpl/outlet.gtpl")
		if err != nil {
			fmt.Println("Error parsing template")
		}
		tpl.Execute(w, d)
	} else {
		http.Redirect(w, r, "/", 302)
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	if IsLoggedIn(r) {
		http.Redirect(w, r, HomeUrl(Users, Username(r)), 302)
		return
	}

	tmpl := template.Must(template.ParseFiles("tmpl/login.gtpl"))
	tmpl.Execute(w, nil)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	session, err := AppSess.Get(r, "session")
	tpl, _ := template.ParseFiles("tmpl/login.gtpl")
	if err != nil {
		tpl.Execute(w, nil)
	} else {
		isLoggedIn := session.Values["loggedIn"]
		if isLoggedIn != true {
			if r.Method == "POST" {
				if ValidUser(r.FormValue("username"), r.FormValue("password"), Users) {
					session.Values["loggedIn"] = true
					session.Values["username"] = r.FormValue("username")
					session.Save(r, w)
					http.Redirect(w, r, HomeUrl(Users, strings.ToUpper(r.FormValue("username"))), 302)
					return
				} else {
					http.Redirect(w, r, "/login", 302)
				}
			} else if r.Method == "GET" {
				tpl.Execute(w, nil)
			}
		} else {
			http.Redirect(w, r, "/", 302)
		}
	}
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	session, err := AppSess.Get(r, "session")
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	session.Options.MaxAge = -1
	session.Values["loggedIn"] = false
	session.Values["username"] = ""
	err = session.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	http.Redirect(w, r, "/", 302)
}

func main() {
	var err error

	flag.Parse()
	Users, err = loadUsers("./config/users.json")
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

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
	mux.HandleFunc("/operator", RequiresLogin(operatorHandler))
	mux.HandleFunc("/dashboard", RequiresLogin(dashboardHandler))
	mux.HandleFunc("/outlet", RequiresLogin(outletHandler))
	mux.HandleFunc("/terminal", RequiresLogin(terminalHandler))
	mux.HandleFunc("/login", loginHandler)
	mux.HandleFunc("/logout", logoutHandler)
	mux.HandleFunc("/", homeHandler)


	// websocket route
	mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})

	if err := http.ListenAndServe(*addr, mux); err != nil {
		log.Fatal("Error starting: ", err)
	}
}