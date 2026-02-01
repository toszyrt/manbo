package openclaw

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/gorilla/websocket"
)

func NewWsClient(config *OpenclawConfig) *WSClient {
	c := &WSClient{
		Config: config,
	}
	if c.Config == nil {
		exePath, err := os.Executable()
		if err != nil {
			log.Fatal(err)
		}
		dir := filepath.Dir(exePath)
		tokenPath := filepath.Join(dir, "openclaw.json")
		data, err := os.ReadFile(tokenPath)
		if err != nil {
			log.Fatal(err)
		}
		var cfg OpenclawConfig
		if err := json.Unmarshal(data, &cfg); err != nil {
			log.Fatal(err)
		}
		c.Config = &cfg
	}
	if c.Config.Ip == "" || c.Config.Port == "" || c.Config.Token == "" {
		log.Fatal("OpenclawConfig has empty fields")
	}
	return c
}

func (c *WSClient) Read() (Message, error) {
	_, data, err := c.conn.ReadMessage()
	if err != nil {
		return Message{}, err
	}
	// log.Printf("[DEBUG] WSClient ⬇ recv: %s\n", data)
	var msg Message
	if err := json.Unmarshal(data, &msg); err != nil {
		return Message{}, err
	}
	return msg, nil
}

func (c *WSClient) Send(msg Message) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	// log.Printf("[DEBUG] WSClient ⬆ send: %s\n", data)
	return c.conn.WriteMessage(websocket.TextMessage, data)
}

func (c *WSClient) ConnectAndHandshake() {
	// 1. 连接
	url := fmt.Sprintf("ws://%s:%s", c.Config.Ip, c.Config.Port)
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Fatal(err)
	}
	c.conn = conn
	msg, err := c.Read()
	if err != nil {
		log.Fatal(err)
	}
	if msg.Type != "event" || msg.Event != "connect.challenge" {
		log.Fatal("invalid msg")
	}
	// 2. 握手
	if err := c.Send(connectReq(c.Config.Token)); err != nil {
		log.Fatal(err)
	}
	msg, err = c.Read()
	if err != nil {
		log.Fatal(err)
	}
	if msg.Type != "res" || msg.ID != "connect" || !msg.Ok {
		log.Fatal("invalid msg")
	}
}

func (c *WSClient) SendQueryAndWait(query string) (reply string, err error) {
	defer func() {
		if err != nil {
			log.Printf("❌ sendQueryAndWait failed: %v", err)
		} else {
			log.Println("🗨 query:  ", query)
			log.Println("🤖 reply: ", reply)
		}
	}()

	locked := make(chan struct{})
	go func() {
		c.mu.Lock()
		close(locked)
	}()
	select {
	case <-locked:
	case <-time.After(5 * time.Second):
		return "", errors.New("get lock timeout")
	}
	defer c.mu.Unlock()

	timeout := time.After(30 * time.Second)
	resultCh := make(chan struct {
		text string
		err  error
	}, 1)

	go func() {
		text, err := c.sendQueryAndWait(query)
		resultCh <- struct {
			text string
			err  error
		}{text, err}
	}()

	select {
	case r := <-resultCh:
		return r.text, r.err
	case <-timeout:
		return "", errors.New("wait timeout")
	}
}

func (c *WSClient) sendQueryAndWait(query string) (string, error) {
	if err := c.Send(agentReq(query)); err != nil {
		return "", err
	}
	var lastText string

	for {
		msg, err := c.Read()
		if err != nil {
			return "", err
		}
		if msg.Type != "event" {
			continue
		}
		event := msg.Event
		payload := msg.Payload
		var payloadMap map[string]any
		if err := json.Unmarshal([]byte(payload), &payloadMap); err != nil {
			return "", err
		}
		switch event {
		case "agent":
			stream, _ := payloadMap["stream"].(string)
			switch stream {
			case "assistant":
				if data, ok := payloadMap["data"].(map[string]any); ok {
					if text, ok := data["text"].(string); ok {
						lastText = text
					}
				}
			case "lifecycle":
				if data, ok := payloadMap["data"].(map[string]any); ok {
					phase, _ := data["phase"].(string)
					if phase == "end" {
						return lastText, nil
					}
					if phase == "error" {
						return "", errors.New("lifecycle error")
					}
				}
			}
		}
	}
}
