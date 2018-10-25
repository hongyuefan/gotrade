package types

type RspVerifyCodeGen struct {
	VerifyCodePng string `json:"verifycode_png"`
	VerifyCodeId  string `json:"verifycode_id"`
}
