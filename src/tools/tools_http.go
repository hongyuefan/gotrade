package tools

import (
	"bytes"
	"net/http"
)

func Post(data []byte, url string) (body *http.Response, err error) {

	client := &http.Client{}

	buff := bytes.NewBuffer(data)

	req, err := http.NewRequest("POST", url, buff)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Connection", "close")

	return client.Do(req)
}
func Get(url string) (body *http.Response, err error) {
	return http.Get(url)
}

func UploadFile(fileName string, fileData []byte, url string) (*http.Response, error) {

	rqbody := new(bytes.Buffer)

	mWriter := multipart.NewWriter(rqbody)

	iow, err := mWriter.CreateFormFile("file", fileName)
	if err != nil {
		return nil, err
	}

	iow.Write(fileData)

	mWriter.Close()

	req, err := http.NewRequest("POST", url, rqbody)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", mWriter.FormDataContentType())

	client := &http.Client{}

	return client.Do(req)
}
