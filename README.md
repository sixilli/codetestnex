## Code Test
All of the code exists within main.go, db.go and routes.go. The extra file hub.go was part of implementing websockets that I didn't have the time to get around to. The websockets would have been used for the endpoint /live that would return live updates. The rest of the CRUD operations should work using BuntDB as the in memory database.

The endpoints only take certain types of requests

/         GET

/create   POST

/read     GET

/update   PUT

/delete   ID
