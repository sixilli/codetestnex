package main

import (
	"log"
	"net/http"
    "html/template"
    "encoding/json"
	"time"

	"github.com/gorilla/websocket"
)

// Almost all this code is taken from gorilla's documentation
// Slightly altered to work with my in memory db implementation

const (
	// Time allowed to write the file to the client.
    writeWait = 10 * time.Second

    // Time allowed to read the next pong message from the client.
    pongWait = 60 * time.Second

	// Send pings to client with this period. Must be less than pongWait.
    pingPeriod = (pongWait * 9) / 10

    // Poll file for changes with this period.
    dbPeroid = 10 * time.Second
)

var (
    liveTempl = template.Must(template.New("").Parse(liveHTML))
    upgrader = websocket.Upgrader{
        ReadBufferSize: 1024,
        WriteBufferSize: 1024,
    }
)

func getDbChanges(lastTime string, dbm DBMem) ([]byte, error) {
    if !(dbm.History[len(dbm.History)-1].Timestamp == lastTime) {
        newChange := dbm.History[len(dbm.History)-1]
        bytes, err := json.Marshal(newChange) 
        if err != nil{
            log.Println(err)
        }
        return bytes, nil
    }

    return nil, nil
}

func reader(ws *websocket.Conn) {
	defer ws.Close()
	ws.SetReadLimit(512)
	ws.SetReadDeadline(time.Now().Add(pongWait))
	ws.SetPongHandler(func(string) error { ws.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, _, err := ws.ReadMessage()
		if err != nil {
			break
		}
	}
}

func writer(ws *websocket.Conn, lastTime string) {
	lastError := ""
	pingTicker := time.NewTicker(pingPeriod)
	dbTicker := time.NewTicker(dbPeroid)
	defer func() {
		pingTicker.Stop()
		dbTicker.Stop()
		ws.Close()
	}()
	for {
		select {
		case <-dbTicker.C:
			var p []byte
			var err error

            changes, err := getDbChanges(lastTime, dm)

			if err != nil {
				if s := err.Error(); s != lastError {
					lastError = s
					p = []byte(lastError)
				}
			} else {
				lastError = ""
			}

			if p != nil {
				ws.SetWriteDeadline(time.Now().Add(writeWait))
				if err := ws.WriteMessage(websocket.TextMessage, changes); err != nil {
					return
				}
			}
		case <-pingTicker.C:
			ws.SetWriteDeadline(time.Now().Add(writeWait))
			if err := ws.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}

func ServeWs(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		if _, ok := err.(websocket.HandshakeError); !ok {
			log.Println(err)
		}
		return
	}

    lastMod := time.Now().String()
	//var lastChange time.Time
	//if n, err := strconv.ParseInt(r.FormValue("lastMod"), 16, 64); err == nil {
		//lastMod = time.Unix(0, n)
	//}
	go writer(ws, lastMod)
	reader(ws)
}

// ServeLive websocket endpoint for live updates
func ServeLive(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	data, err := getDbChanges(time.Now().String(), dm)
	if err != nil {
        log.Println(err)
	}
    // Some magic numbers, better to return a struct from getDbChanges
	var v = struct {
		Host       string
		Data       string
		LastMod    string
	}{
		r.Host,
        string(data),
		string(data),
	}
	liveTempl.Execute(w, &v)
}

const liveHTML = `<!DOCTYPE html>
<html lang="en">
    <head>
        <title>WebSocket Example</title>
    </head>
    <body>
        <pre id="fileData">{{.Data}}</pre>
        <script type="text/javascript">
            (function() {
                var data = document.getElementById("fileData");
                var conn = new WebSocket("ws://{{.Host}}/ws?lastMod={{.LastMod}}");
                conn.onclose = function(evt) {
                    data.textContent = 'Connection closed';
                }
                conn.onmessage = function(evt) {
                    console.log('file updated');
                    data.textContent = evt.data;
                }
            })();
        </script>
    </body>
</html>`
