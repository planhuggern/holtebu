package osloconquest

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWebSocketHandlerUpgradesWebSocketRequest(t *testing.T) {
	server := httptest.NewServer(WebSocketHandler())
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
