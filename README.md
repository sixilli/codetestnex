## Code Test
The API can handle all crud operations and there's an endpoint that updates live with changes to the database after connection. 

| File        | Description     |
| ------------- |-------------:|
| dbMem.go      | Implementation of in memory database      |
| routes.go | Handles all routing for the API      |
| client.go | Websockt client      |
| hub.go | Responsible for keeping track of websocket clients      |
| ws.go | Websocket endpoints and reader/writer functions      |

/         GET

/create   POST

/read     GET

/update   PUT

/delete   ID

/live
