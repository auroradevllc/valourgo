package signalr

import "encoding/json"

var (
	jsonSerializer = new(JSONSerializer)
)

type Serializer interface {
	Serialize(value any) ([]byte, error)
}

type JSONSerializer struct {
}

func (s *JSONSerializer) Serialize(value any) ([]byte, error) {
	return json.Marshal(value)
}
