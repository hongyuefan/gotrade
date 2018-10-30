package jsonprocess

import (
	"encoding/json"
)

type JsonProcess struct {
}

func NewJsonProcess() *JsonProcess {
	return &JsonProcess{}
}

func (j *JsonProcess) UnMarshal(data []byte, msg interface{}) error {
	return json.Unmarshal(data, msg)
}

func (j *JsonProcess) Marshal(msg interface{}) ([]byte, error) {
	return json.Marshal(msg)
}
