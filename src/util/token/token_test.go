package token

import (
	"testing"
)

func TestTokenGenerate(t *testing.T) {

	token, _ := TokenGenerate(1, 1)

	t.Log(token)

	userId, _ := TokenValidate(token)

	t.Log(userId)
}
