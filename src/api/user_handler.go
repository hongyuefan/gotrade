package api

import (
	"crypto/md5"
	"encoding/hex"
	models "models/common"
	"time"
	"types"
	"util/token"
	vc "util/verifycode"

	"github.com/gin-gonic/gin"
)

func (h *Handlers) HandlerVerifyCode(c *gin.Context) {
	var (
		rspVcode         *types.RspVerifyCodeGen
		capId, base64Png string
	)

	capId, base64Png = vc.CodeGenerate(60, 240, 4)

	rspVcode = &types.RspVerifyCodeGen{
		VerifyCodePng: base64Png,
		VerifyCodeId:  capId,
	}

	HandleSuccessMsg(c, "HandlerVerifyCode", rspVcode)
}
func (h *Handlers) HandlerUpdatePayPass(c *gin.Context) {
	var (
		err       error
		reqUpdate *types.ReqUpdatePassword
		user      *models.User
		userId    int64
	)
	if err = c.BindJSON(reqUpdate); err != nil {
		goto errDeal
	}
	if userId, err = token.TokenValidate(reqUpdate.Token); err != nil {
		goto errDeal
	}

	if user, err = models.GetUserById(userId); err != nil {
		goto errDeal
	}
	if user.PayPass == "" {
		if user.Password != MD5(reqUpdate.OldPass) {
			err = types.Error_Password_Wrong
			goto errDeal
		}
	} else {
		if user.PayPass != MD5(reqUpdate.OldPass) {
			err = types.Error_Password_Wrong
			goto errDeal
		}
	}
	user = &models.User{
		Id:      userId,
		PayPass: MD5(reqUpdate.NewPass),
	}
	HandleSuccessMsg(c, "HandlerUpdatePayPass", "success")
	return
errDeal:
	HandleErrorMsg(c, "HandlerUpdatePayPass", err)
	return
}

func (h *Handlers) HandlerUpdatePassword(c *gin.Context) {
	var (
		err       error
		reqUpdate *types.ReqUpdatePassword
		user      *models.User
		userId    int64
	)
	if err = c.BindJSON(reqUpdate); err != nil {
		goto errDeal
	}
	if userId, err = token.TokenValidate(reqUpdate.Token); err != nil {
		goto errDeal
	}

	if user, err = models.GetUserById(userId); err != nil {
		goto errDeal
	}
	if MD5(reqUpdate.OldPass) != user.Password {
		err = types.Error_Password_Wrong
		goto errDeal
	}
	user = &models.User{
		Id:       userId,
		Password: MD5(reqUpdate.NewPass),
	}

	if err = models.UpdateUserById(user, "password"); err != nil {
		goto errDeal
	}
	HandleSuccessMsg(c, "HandlerUpdatePasswoed", "success")
	return
errDeal:
	HandleErrorMsg(c, "HandlerUpdatePasswoed", err)
	return
}

func (h *Handlers) HandlerUserLogin(c *gin.Context) {
	var (
		err    error
		reqReg types.ReqUserRegAndLogin
		rspReg *types.RspUserRegAndLogin
		user   *models.User
		tok    string
	)

	if err = c.BindJSON(&reqReg); err != nil {
		goto errDeal
	}

	user = &models.User{
		Username: reqReg.UserName,
		Password: MD5(reqReg.PassWord),
		IsUsed:   true,
	}

	if err = models.GetUser(user, "username", "password"); err != nil {
		goto errDeal
	}

	if tok, err = token.TokenGenerate(user.Id, h.ExpireTimeToken); err != nil {
		goto errDeal
	}

	rspReg = &types.RspUserRegAndLogin{Token: tok}

	HandleSuccessMsg(c, "HandlerUserLogin", rspReg)

	return
errDeal:
	HandleErrorMsg(c, "HandlerUserLogin", err)
	return
}

func (h *Handlers) HandlerUserRegist(c *gin.Context) {
	var (
		err           error
		reqReg        types.ReqUserRegAndLogin
		rspReg        *types.RspUserRegAndLogin
		appId, appKey string
		user          *models.User
		uID           int64
		tok           string
	)

	if err = c.BindJSON(&reqReg); err != nil {
		goto errDeal
	}

	if !vc.CodeValidate(reqReg.VerifyCodeId, reqReg.VerifyCode) {
		err = types.Error_Verifycode_Wrong
		goto errDeal
	}

	user = &models.User{
		Username:   reqReg.UserName,
		Password:   MD5(reqReg.PassWord),
		AppId:      appId,
		AppKey:     appKey,
		Createtime: time.Now().Unix(),
		IsUsed:     true,
	}

	if uID, err = models.AddUser(user); err != nil {
		goto errDeal
	}

	if tok, err = token.TokenGenerate(uID, h.ExpireTimeToken); err != nil {
		goto errDeal
	}

	rspReg = &types.RspUserRegAndLogin{Token: tok}

	HandleSuccessMsg(c, "HandlerUserRegist", rspReg)

	return
errDeal:
	HandleErrorMsg(c, "HandlerUserRegist", err)
	return
}

func MD5(input string) string {
	output := md5.Sum([]byte(input))
	return hex.EncodeToString(output[:])
}
