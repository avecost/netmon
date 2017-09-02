package main

import (
	"time"
	"encoding/json"
)

// Uptime will contain event name 'e', and current time 't'
type Uptime struct {
	Event string
	ServerT string
}

// PushTime return the struct of current time with event
func PushTime() (t []byte) {
	tmp := Uptime{Event: "TIME-UPDATE", ServerT: time.Now().Local().Format("Mon Jan 2 2006 03:04:05 PM")}
	t, err := json.Marshal(tmp)
	if err != nil {
		return nil
	}
	return t
}