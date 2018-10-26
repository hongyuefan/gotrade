package sign

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"io"
)

func HMacSha256(msg string, key []byte) string {

	h := hmac.New(sha256.New, key)

	io.WriteString(h, msg)

	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}
