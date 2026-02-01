package openclaw

import (
	"encoding/json"
	"sync"

	"github.com/gorilla/websocket"
)

type WSClient struct {
	Config *OpenclawConfig
	conn   *websocket.Conn
	mu     sync.Mutex
}

type OpenclawConfig struct {
	Ip    string `json:"ip,omitempty"`
	Port  string `json:"port,omitempty"`
	Token string `json:"token,omitempty"`
}

type Message struct {
	Type   string          `json:"type"`
	ID     string          `json:"id,omitempty"`
	Event  string          `json:"event,omitempty"`
	Method string          `json:"method,omitempty"`
	Params json.RawMessage `json:"params,omitempty"`
	Ok     bool            `json:"ok,omitempty"`
	Error  *struct {
		Message string `json:"message"`
	} `json:"error,omitempty"`
	Payload json.RawMessage `json:"payload,omitempty"`
}
