package jsonprocess

import (
	"encoding/json"
	"util/wclient"
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

func (j *JsonProcess) Route(data interface{}, agent wclient.Agent) error {
	return nil
}
