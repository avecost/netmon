package main

import (
	"log"
	"time"
	"encoding/json"
	"fmt"
)

// Hub contains the information for
// websocket channels/clients
type Hub struct {
	// Dashboard data
	iest []Franchisee

	// registered clients
	clients map[*Client]bool

	// Rooms for client connections.
	rooms map[string]map[*Client]bool

	// Channel for rooms
	bcroom chan []byte

	// Dashboard
	dashboard chan []byte

	// Inbound messages from the clients
	broadcast chan []byte

	// Register requests from the clients
	register chan *Client

	// Unregister request from the clients
	unregister chan *Client
}

// NetmonHeader is the packet content expected
// from netmon-c client
type NetmonHeader struct {
	Event  string	 	`json:"event"`
	Outlet string 		`json:"outlet"`
	Acct   string 		`json:"acct"`
	Privip string 		`json:"privateip"`
	Pubip  string 		`json:"publicip"`
	Os     string 		`json:"os"`
}

type iestStat struct {
	Outlet		int		`json:"Outlet"`
	Terminal 	int		`json:"Terminal"`
	Online 		int		`json:"Online"`
}

type Outlet struct {
	Name		string	`json:"Name"`
	Terminal 	int		`json:"Terminal"`
	Online 		int		`json:"Online"`
}

func newHub() *Hub {
	iestJson, err := loadOutlet(APP_DIR + "config/outlet.json")
	if err != nil {
		log.Printf("Error: %s\n", err)
		panic("error loading config file.")
	}
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
		rooms: 		make(map[string]map[*Client]bool),
		bcroom:		make(chan []byte),
		dashboard:	make(chan []byte),
		iest: 		iestJson,
	}
}

func (h *Hub) removeClient(client *Client) {
	if _, ok := h.clients[client]; ok {
		delete(h.clients, client)
		close(client.send)
	}
	for k, r := range h.rooms {
		if len(r) > 0 {
			if _, ok := h.rooms[k][client]; ok {
				delete(h.rooms[k],client)
			}
		}
	}
	for k, r := range h.rooms {
		fmt.Printf("R: %v c: %d\n", k, len(r))
	}
}

func (h *Hub) run() {
	// ticker check inactive terminal every 15s
	tc := time.NewTicker(time.Second * TICKER_ONLINE_TIME)
	// heartbeat ticker runs every second to
	// correspond to server time change
	hb := time.NewTicker(time.Second * TICKER_SERVER_TIME)
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			h.removeClient(client)
			//if _, ok := h.clients[client]; ok {
			//	delete(h.clients, client)
			//	close(client.send)
			//}
			//for i := range h.rooms {
			//	c := h.rooms[i]
			//	if c == nil {
			//		continue
			//	}
			//	if len(h.rooms[i]) == 0 {
			//		continue
			//	}
			//	if _, ok := h.rooms[i][client]; ok {
			//		delete(h.rooms[i], client)
			//		delete(h.clients, client)
			//		close(client.send)
			//	}
			//}
			fmt.Printf("unreg: %v\n", h.clients)
		case message := <-h.broadcast:
			//var nmH NetmonHeader
			//json.Unmarshal(message, &nmH)
			//log.Printf("<%s> %s %s %s %s %s\n", nmH.Event, nmH.Outlet, nmH.Acct, nmH.Privip, nmH.Pubip, nmH.Os)

			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
			fmt.Printf("broadcast: %v\n", h.clients)
		case dashboard := <-h.dashboard:
			for client := range h.clients {
				select {
				case client.send <- dashboard:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
			fmt.Printf("dashboard: %v\n", h.clients)
		case t := <-h.bcroom:
			for i := range h.rooms {
				c := h.rooms[i]
				if c == nil {
					continue
				}
				if len(h.rooms[i]) == 0 {
					continue
				}
				for u := range h.rooms[i] {
					fmt.Printf("Room: %v\n", h.rooms[i])
					fmt.Printf("bcroom: %v\n", u)

					select {
					case u.send <- t:

						//default:
						//	close(u.send)
						//	delete(h.rooms[i], u)
						//	delete(h.clients, u)
					}
				}
			}
		case <-hb.C:
			t := PushTime()

			for client := range h.clients {
				select {
				case client.send <- t:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
			fmt.Printf("ticker: %v\n", h.clients)
		case <-tc.C:
			go h.chkOnline()
		}
	}
}

func (h *Hub) update(t string) {
	tNow, _ := time.LoadLocation(TIME_ZONE)
	t2 := time.Now().In(tNow).Format(TIME_FORMAT)
	for i := range h.iest {
		for j := range h.iest[i].Outlets {
			ts := h.iest[i].Outlets[j].Terminals
			for k := range ts {
				if ts[k].Account == t {
					ts[k].Status = 1
					if len(ts[k].Online) == 0 {
						ts[k].Online = t2
					}
					ts[k].Lastupdate = t2
				}
			}
		}
	}
}

func (h *Hub) netsum(m map[string]iestStat) {
	var o, t, ol int
	for i := range h.iest {
		opName := h.iest[i].Operator					// operator name
		for j := range h.iest[i].Outlets {
			o = len(h.iest[i].Outlets)					// outlet count
			ts := h.iest[i].Outlets[j].Terminals

			t = t + len(ts)		// terminal count
			for k := range ts {
				if ts[k].Status != 2 {					// we only need online terminal
					ol = ol + 1							// terminal online count
				}
			}
		}
		m[opName] = iestStat{Outlet: o, Terminal: t, Online: ol}
		o, t, ol = 0, 0, 0
	}
}

func (h *Hub) chkOnline() {
	tLoc, _ := time.LoadLocation(TIME_ZONE)
	t2 := time.Now().In(tLoc)

	for i := range h.iest {
		for j := range h.iest[i].Outlets {
			ts := h.iest[i].Outlets[j].Terminals
			for k := range ts {
				if len(ts[k].Lastupdate) > 0 {
					t, _ := time.ParseInLocation(TIME_FORMAT, ts[k].Lastupdate, tLoc)
					if t2.Sub(t) > (time.Minute * CLIENT_TTL) {
						log.Printf("Offline: %s %s\n", h.iest[i].Operator, ts[k].Account)

						ts[k].Status = 2
						ts[k].Online = ""
						ts[k].Lastupdate = ""
					}
				}
			}
		}
	}

	netsum := make(map[string]iestStat)
	h.netsum(netsum)
	dashboard, _ := json.Marshal(&DBHeader{Event: "DB-UPDATE", Netsum: netsum})
	h.dashboard <- dashboard
}

func (h *Hub) createRoom() {
	for i := range h.iest {
		h.rooms[h.iest[i].Operator] = nil
	}
}

func (h *Hub) byOperator() {
	var n string
	var t, ol int

	outletSummary := make(map[string][]*Outlet)
	for i := range h.iest {
		r := h.iest[i].Operator
		for j := range h.iest[i].Outlets {
			n = h.iest[i].Outlets[j].Name
			ts := h.iest[i].Outlets[j].Terminals
			t = t + len(ts)		// terminal count
			for k := range ts {
				if ts[k].Status != 2 {					// we only need online terminal
					ol = ol + 1							// terminal online count
				}
			}
			outletInfo := &Outlet{Name: n, Terminal: t, Online: ol}
			outletSummary[r] = append(outletSummary[r], outletInfo)

			//fmt.Printf("%s %s %d %d\n", r, n, t, ol)
			n, t, ol = "", 0, 0
		}
	}

	outletDB, _ := json.Marshal(&OutletHeader{Event: "OUTLET-UPDATE", Outsum: outletSummary})
	h.bcroom <-outletDB
}