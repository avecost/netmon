package main

import (
	"log"
	"time"
	"encoding/json"
)

// Hub contains the information for
// websocket channels/clients
type Hub struct {
	// registered clients
	clients map[*Client]bool

	// Inbound messages from the clients
	broadcast chan []byte

	// Register requests from the clients
	register chan *Client

	// Unregister request from the clients
	unregister chan *Client

	// Dashboard
	dashboard chan []byte

	// Dashboard data
	iest []Franchisee
}

// NetmonHeader is the packet content expected
// from netmon-c client
type NetmonHeader struct {
	Event  string `json:"event"`
	Outlet string `json:"outlet"`
	Acct   string `json:"acct"`
	Privip string `json:"privateip"`
	Pubip  string `json:"publicip"`
	Os     string `json:"os"`
}

type iestStat struct {
	Outlet		int
	Terminal 	int
	Online 		int
}

func newHub() *Hub {
	iestJson, _ := loadOutlet("./config/outlet.json")

	//iestJson, _ := loadOutlet("/home/whiskie/netmon/config/outlet.json")
	//iestJson, err := loadOutlet("/config/outlet.json")
	//if err != nil {
	//	log.Printf("Error: %s\n", err)
	//	panic("error loading config file.")
	//}
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
		dashboard:	make(chan []byte),
		iest: 		iestJson,
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
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
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
		case dashboard := <-h.dashboard:
			for client := range h.clients {
				select {
				case client.send <- dashboard:
				default:
					close(client.send)
					delete(h.clients, client)
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
