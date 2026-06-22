package server

import "encoding/json"

type Packet struct {
	Event string          `json:"event"`
	Data  json.RawMessage `json:"data,omitempty"`
	AckID string          `json:"ackId,omitempty"`
}

type HandleFunc func(c *Client, data json.RawMessage)
