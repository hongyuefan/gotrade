package app

import (
	"api"

	"fmt"
	"net/http"
	"server/jsonprocess"
	"server/restful"
	"strconv"
	"time"
	"util/sign"

	"util/config"
	"util/log"

	"github.com/astaxie/beego/orm"
	gin "github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

const MasterName = "gotrade"

type ConfigData struct {
	Port              string
	Idls              float64
	LogDir            string
	AppId             string
	AppKey            string
	MchId             string
	IP                string
	NotifyUrl         string
	ExpireTimeToken   float64
	ExpireTimeSession float64
}

type App struct {
	handlers   *api.Handlers
	restServer *restful.RestServer
	closeSig   chan bool
}

var g_ConfigData *ConfigData

func OnInitFlag(c *config.Config) (err error) {

	g_ConfigData = new(ConfigData)
	g_ConfigData.Port = c.GetString("port")
	g_ConfigData.Idls = c.GetFloat("idls")
	g_ConfigData.LogDir = c.GetString("logdir")

	g_ConfigData.AppId = c.GetString("appid")
	g_ConfigData.AppKey = c.GetString("appkey")
	g_ConfigData.MchId = c.GetString("mchid")
	g_ConfigData.IP = c.GetString("ip")
	g_ConfigData.NotifyUrl = c.GetString("notifyUrl")
	g_ConfigData.ExpireTimeToken = c.GetFloat("expire_token")
	g_ConfigData.ExpireTimeSession = c.GetFloat("expire_session")

	if "" == g_ConfigData.Port || 0 == g_ConfigData.Idls || "" == g_ConfigData.LogDir {
		return fmt.Errorf("config not right")
	}
	return

}

func (app *App) RegistRestServer() {

	app.restServer = restful.NewRestServer("https://www.okex.com", jsonprocess.NewJsonProcess())

	app.restServer.RegistInterface("order", "/api/futures/v3/order", restful.Method_Post)
	app.restServer.RegistInterface("position", "/api/futures/v3/position", restful.Method_Get)
	app.restServer.RegistInterface("sigleposition", "/api/futures/v3/%v/position", restful.Method_Get)

}
func (app *App) CotrolHandlers() {

	app.RegistRestServer()

	mm := make(map[string]interface{}, 0)
	heards := make(map[string]string, 0)

	var unixMic int = int(time.Now().UnixNano() / 1000000)
	timeStamp := strconv.Itoa(unixMic)
	timeStamp = timeStamp[:len(timeStamp)-3] + "." + timeStamp[len(timeStamp)-3:]

	heards["OK-ACCESS-KEY"] = "342d1884-db81-4a9c-8535-1d4351965adf"
	heards["OK-ACCESS-SIGN"] = sign.HMacSha256(timeStamp+"GET"+"/api/futures/v3/position", []byte("3628818392EC421EF456070057E0F9CF"))
	heards["OK-ACCESS-TIMESTAMP"] = timeStamp
	heards["OK-ACCESS-PASSPHRASE"] = "IMDANDAN"
	heards["contentType"] = "application/json"

	body, err := app.restServer.SynCall("position", heards, mm)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(body))
	//	app.closeSig = make(chan bool, 1)

	//	sagent := okex.NewAgentLogin(256)

	//	klGate := wshb.NewGate("wss://real.okex.com:10440/websocket/okexapi", 1, 1024, 65536, 5*time.Second, 5*time.Second, true, sagent)

	//	go klGate.Run(app.closeSig)

	//	var lg mo.ReqFurtureLogin
	//	var unixMic int = int(time.Now().UnixNano() / 1000000)
	//	timeStamp := strconv.Itoa(unixMic)
	//	timeStamp = timeStamp[:len(timeStamp)-3] + "." + timeStamp[len(timeStamp)-3:]
	//	fmt.Println(timeStamp)
	//	lg.Event = "login"
	//	lg.Params.ApiKey = "342d1884-db81-4a9c-8535-1d4351965adf"
	//	lg.Params.PassPhrase = "IMDANDAN"
	//	lg.Params.Sign = sign.HMacSha256(timeStamp+"GET"+"/users/self/verify", []byte("3628818392EC421EF456070057E0F9CF"))
	//	lg.Params.TimeStamp = timeStamp

	//	sagent.WriteMsg(lg)
}

func (app *App) OnStart(c *config.Config) error {

	if err := OnInitFlag(c); err != nil {
		return err
	}

	if _, err := log.NewLog(g_ConfigData.LogDir, MasterName, 0); err != nil {
		return err
	}

	orm.RegisterDataBase("default", "mysql", "root:@tcp(127.0.0.1:3306)/payplate")

	app.handlers = api.NewHandlers(g_ConfigData.AppId, g_ConfigData.AppKey, g_ConfigData.MchId, g_ConfigData.IP, g_ConfigData.NotifyUrl, int64(g_ConfigData.ExpireTimeSession), int64(g_ConfigData.ExpireTimeToken))

	router := gin.Default()

	v0 := router.Group("/v0")
	{
		v0.GET("/health", app.handlers.HandlerGet)
	}

	v1 := router.Group("/v1")
	{
		v1.POST("/post", app.handlers.HandlerPost)
		v1.GET("/get", app.handlers.HandlerGet)

		v1.POST("/payplat/user/regist", app.handlers.HandlerUserRegist)
		v1.GET("/payplat/verifycode", app.handlers.HandlerVerifyCode)
	}

	fmt.Println("Listen:", g_ConfigData.Port)

	app.CotrolHandlers()

	http.ListenAndServe(":"+g_ConfigData.Port, router)

	return nil
}

func (app *App) Shutdown() {
	app.handlers.OnClose()
	close(app.closeSig)
	fmt.Println("server shutdown")
}
