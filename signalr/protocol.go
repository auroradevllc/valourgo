package signalr

import (
	"encoding/json"
)

const (
	messageTypeInvocation = 1
	messageTypeCompletion = 3
	messageTypePing       = 6
)

type Frame struct {
	Type           int               `json:"type"`
	InvocationID   string            `json:"invocationId,omitempty"`
	Target         string            `json:"target,omitempty"`
	Arguments      []json.RawMessage `json:"arguments,omitempty"`
	Result         json.RawMessage   `json:"result,omitempty"`
	Error          string            `json:"error,omitempty"`
	AllowReconnect bool              `json:"allowReconnect,omitempty"`
}

func parseFrames(data []byte) ([]Frame, error) {
	var frames []Frame
	for _, part := range split(data) {
		var f Frame

		if err := json.Unmarshal(part, &f); err != nil {
			return nil, err
		}

		frames = append(frames, f)
	}
	return frames, nil
}

// SignalR uses ASCII 0x1e as a record separator
func split(b []byte) [][]byte {
	var out [][]byte
	start := 0
	for i, c := range b {
		if c == 0x1e {
			out = append(out, b[start:i])
			start = i + 1
		}
	}
	return out
}
