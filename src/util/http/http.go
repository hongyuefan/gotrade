package http

import (
	"bytes"
	"fmt"
	"io"
	hp "net/http"
)

/*
contentType:application/json
*/
func Post(url string, headers map[string]string, data []byte) ([]byte, error) {

	client := &hp.Client{}

	buff := bytes.NewBuffer(data)

	req, err := hp.NewRequest("POST", url, buff)
	if err != nil {
		return nil, err
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	rsp, err := client.Do(req)
	if rsp.StatusCode != hp.StatusOK {
		return nil, fmt.Errorf("Post %v Error code:%v,msg:%v", url, rsp.StatusCode, rsp.Status)
	}

	buf := new(bytes.Buffer)

	io.Copy(buf, rsp.Body)

	return buf.Bytes(), nil
}

func Get(url string, headers map[string]string) ([]byte, error) {

	client := &hp.Client{}

	req, err := hp.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	rsp, err := client.Do(req)
	if rsp.StatusCode != hp.StatusOK {
		return nil, fmt.Errorf("GET %v Error code:%v,msg:%v", url, rsp.StatusCode, rsp.Status)
	}

	buf := new(bytes.Buffer)

	io.Copy(buf, rsp.Body)

	return buf.Bytes(), nil
}
