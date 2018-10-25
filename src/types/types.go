package types

import (
	"encoding/json"
	"errors"
)

type RspCommon struct {
	Success bool            `json:"success"`
	Message string          `json:"message"`
	Data    json.RawMessage `json:"data"`
}

var (
	Error_Verifycode_Wrong = errors.New("验证码错误")
	Error_Password_Wrong   = errors.New("密码错误")
)
