package api

import (
	"encoding/json"
	"fmt"
	"types"
	"util/log"
	"util/session"

	"github.com/gin-gonic/gin"
)

type Handlers struct {
	AppId           string
	AppKey          string
	MchId           string
	NotifyUrl       string
	IP              string
	Session         *session.Session
	ExpireTimeToken int64
}

func NewHandlers(appId, appKey, mchId, IP, notifyUrl string, expireTimeSession, expireTimeToken int64) *Handlers {

	return &Handlers{
		AppId:           appId,
		AppKey:          appKey,
		MchId:           mchId,
		IP:              IP,
		NotifyUrl:       notifyUrl,
		Session:         session.NewSession(expireTimeSession, 5),
		ExpireTimeToken: expireTimeToken,
	}
}

func (h *Handlers) OnClose() {

}

func (h *Handlers) HandlerPost(c *gin.Context) {
	var (
		err error
	)
	if err != nil {
		goto errDeal
	}
	return
errDeal:
	HandleErrorMsg(c, "HandlerPost", err)
	return
}

func (h *Handlers) HandlerGet(c *gin.Context) {
	var (
		err error
	)
	if err != nil {
		goto errDeal
	}
	return
errDeal:
	HandleErrorMsg(c, "HandlerGet", err)
	return
}

func HandleNotifySuccess(c *gin.Context, requestType string, out interface{}) {

	jsonData, _ := json.Marshal(out)

	c.XML(200, &types.RspNotify{
		ReturnCode: "SUCCESS",
		ReturnMsg:  "OK",
	})

	logMsg := fmt.Sprintf("type[%s] From [%s] Params [%s]", requestType, c.Request.RemoteAddr, jsonData)

	log.GetLog().LogDebug(logMsg)

	return
}

func HandleNotifyFailed(c *gin.Context, requestType string, err error) {

	c.XML(200, &types.RspNotify{
		ReturnCode: "FAILED",
		ReturnMsg:  err.Error(),
	})

	logMsg := fmt.Sprintf("type[%s] From [%s] Error [%s] ", requestType, c.Request.RemoteAddr, err.Error())

	log.GetLog().LogError(logMsg)
}

func HandleErrorMsg(c *gin.Context, requestType string, err error) {

	msg := types.RspCommon{
		Success: false,
		Message: err.Error(),
	}

	c.JSON(200, msg)

	logMsg := fmt.Sprintf("type[%s] From [%s] Error [%s] ", requestType, c.Request.RemoteAddr, msg)

	log.GetLog().LogError(logMsg)

	return
}

func HandleSuccessMsg(c *gin.Context, requestType string, out interface{}) {

	jsonData, err := json.Marshal(out)

	if err != nil {
		HandleErrorMsg(c, requestType, err)
		return
	}

	msg := types.RspCommon{
		Success: true,
		Message: "success",
		Data:    jsonData,
	}

	c.JSON(200, msg)

	logMsg := fmt.Sprintf("type[%s] From [%s] Params [%s]", requestType, c.Request.RemoteAddr, jsonData)

	log.GetLog().LogDebug(logMsg)

	return
}
