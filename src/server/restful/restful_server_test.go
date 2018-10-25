package restful

import (
	"server/jsonprocess"
	"testing"
)

func TestBench(t *testing.T) {

	s := NewRestServer(jsonprocess.NewJsonProcess())

	if err := s.RegistInterface("getfutures", "https://www.okex.com/api/futures/v3/instruments", "application/json", Method_Get); err != nil {
		t.Log(err)
		return
	}

	mmap := make(map[string]interface{})

	if rsp, err := s.SynCall("getfutures", mmap); err != nil {
		t.Log(err)
		return
	} else {
		t.Log(string(rsp))
	}

	return
}
