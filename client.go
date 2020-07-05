package main

import (
	"bytes"
    "fmt"
	"log"
	"net/http"
    "strconv"
    "strings"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	hub *Hub

	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	history chan []byte

    // Last update
    lastChange string
}

func lastChange(lastChange string) (dbEvent, bool) {
    lastHistoryItem := dm.History[len(dm.History)-1]
    if strings.Compare(lastHistoryItem.Timestamp, lastChange) != 0 {
        return lastHistoryItem, true
    }
    return lastHistoryItem, false
}

func (c *Client) watchDb() {
    defer func() {
            c.hub.unregister <- c
            c.conn.Close()
    }()
    c.conn.SetReadLimit(maxMessageSize)
    c.conn.SetReadDeadline(time.Now().Add(pongWait))
    c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
    for {
        d, newChange := lastChange(c.lastChange)
        id := strconv.Itoa(d.ID)
        if newChange == true {
            fmt.Println("Found change")
            c.lastChange = d.Timestamp
            message := []byte(d.Op + id + d.Data + d.Timestamp)
            message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
            c.hub.broadcast <- message
        }
    }
}
// writePump pumps messages from the hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.history:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued chat messages to the current websocket message.
			n := len(c.history)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-c.history)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// serveWs handles websocket requests from the peer.
func ServeWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

    startTime := dm.History[len(dm.History)-1].Timestamp
	client := &Client{
        hub: hub, 
        conn: conn, 
        history: make(chan []byte, 256),
        lastChange: startTime,
    }
	client.hub.register <- client

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
    go client.watchDb()
	go client.writePump()
}
