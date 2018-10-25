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
func Post(url string, contentType string, data []byte) ([]byte, error) {

	client := &hp.Client{}

	buff := bytes.NewBuffer(data)

	req, err := hp.NewRequest("POST", url, buff)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", contentType)

	rsp, err := client.Do(req)
	if rsp.StatusCode != hp.StatusOK {
		return nil, fmt.Errorf("Post %v Error code:%v,msg:%v", url, rsp.StatusCode, rsp.Status)
	}

	buf := new(bytes.Buffer)

	io.Copy(buf, rsp.Body)

	return buf.Bytes(), nil
}

func Get(url string, contentType string) ([]byte, error) {
	rsp, err := hp.Get(url)
	if err != nil {
		return nil, err
	}
	if rsp.StatusCode != hp.StatusOK {
		return nil, fmt.Errorf("Get %v Error code:%v,msg:%v", url, rsp.StatusCode, rsp.Status)
	}
	buf := new(bytes.Buffer)

	io.Copy(buf, rsp.Body)

	return buf.Bytes(), nil
}
