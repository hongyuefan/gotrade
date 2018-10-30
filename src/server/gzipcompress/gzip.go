package gzipcompress

import (
	"bytes"
	"compress/flate"
	"io/ioutil"
)

type MsgGZip struct{}

func NewMsgGZip() *MsgGZip {
	return new(MsgGZip)
}

func (gz *MsgGZip) Compress(in []byte) ([]byte, error) {
	return nil, nil
}

func (gz *MsgGZip) UnCompress(in []byte) ([]byte, error) {
	reader := flate.NewReader(bytes.NewReader(in))
	defer reader.Close()

	return ioutil.ReadAll(reader)
}
