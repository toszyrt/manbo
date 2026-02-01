package openclaw

import (
	"encoding/json"

	"github.com/google/uuid"
)

func connectReq(token string) Message {
	return Message{
		Type:   "req",
		ID:     "connect",
		Method: "connect",
		Params: json.RawMessage(`{
			"minProtocol": 3,
			"maxProtocol": 3,
			"client": {
				"id":       "gateway-client",
				"version":  "0.1.0",
				"platform": "macair",
				"mode":     "backend"
			},
			"role":   "operator",
			"scopes": ["operator.read", "operator.write"],
			"auth": {
				"token": "` + token + `"
			},
			"locale":    "zh-CN",
			"userAgent": "gateway-client-go"
		}`),
	}
}

func agentReq(text string) Message {
	return Message{
		Type:   "req",
		ID:     "agent",
		Method: "agent",
		Params: json.RawMessage(`{
			"message": "` + text + `",
			"agentId": "main",
			"sessionKey": "xiaoai-openclaw-session",
			"deliver": false,
			"idempotencyKey": "` + uuid.NewString() + `"
		}`),
	}
}
