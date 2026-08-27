package osloconquest

import (
	"net/http"

	"github.com/coder/websocket"
)

func WebSocketHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := websocket.Accept(w, r, nil)
		if err != nil {
			http.Error(w, "Failed to upgrade to WebSocket", http.StatusInternalServerError)
			return
		}

		// Handle WebSocket connection here (e.g., read/write messages)

		conn.Close(websocket.StatusNormalClosure, "")
	})
}
