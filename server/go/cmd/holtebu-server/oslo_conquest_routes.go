package main

import (
	"net/http"
	"holtebu-server/osloconquest"
)

func registerOsloConquestRoutes(mux *http.ServeMux)  {
	mux.Handle("/ws/oslo-conquest/", osloconquest.WebSocketHandler())


}
