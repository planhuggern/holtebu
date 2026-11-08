package osloconquest

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
)

func TestWebSocketHandlerUpgradesWebSocketRequest(t *testing.T) {
	server := httptest.NewServer(WebSocketHandler(NewRoomStore()))
	defer server.Close()

	request, err := http.NewRequest(http.MethodGet, server.URL, nil)
	if err != nil {
		t.Fatalf("NewRequest returned an error: %v", err)
	}
	request.Header.Set("Connection", "Upgrade")
	request.Header.Set("Upgrade", "websocket")
	request.Header.Set("Sec-WebSocket-Version", "13")
	request.Header.Set("Sec-WebSocket-Key", "dGhlIHNhbXBsZSBub25jZQ==")

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		t.Fatalf("Do returned an error: %v", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusSwitchingProtocols {
		t.Errorf("status code = %d, want %d", response.StatusCode, http.StatusSwitchingProtocols)
	}
}

func TestWebSocketHandlerListsRooms(t *testing.T) {
	server := httptest.NewServer(WebSocketHandler(NewRoomStore()))
	defer server.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	webSocketURL := "ws" + strings.TrimPrefix(server.URL, "http")
	connection, _, err := websocket.Dial(ctx, webSocketURL, nil)
	if err != nil {
		t.Fatalf("Dial returned an error: %v", err)
	}
	defer connection.CloseNow()

	if err := wsjson.Write(ctx, connection, map[string]string{"type": "list_rooms"}); err != nil {
		t.Fatalf("Write returned an error: %v", err)
	}

	var response struct {
		Type  string     `json:"type"`
		Rooms []RoomInfo `json:"rooms"`
	}
	if err := wsjson.Read(ctx, connection, &response); err != nil {
		t.Fatalf("Read returned an error: %v", err)
	}
	if response.Type != "room_list" {
		t.Errorf("response type = %q, want %q", response.Type, "room_list")
	}
	if len(response.Rooms) != 0 {
		t.Errorf("room count = %d, want 0", len(response.Rooms))
	}
}
