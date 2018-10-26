package sign

import (
	"testing"
)

func TestHMacSha256(t *testing.T) {
	t.Log(HMacSha256("123", []byte("3434")))
}
