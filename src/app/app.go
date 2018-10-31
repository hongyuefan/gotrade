package app

import (
	"fmt"
	"net/http"
	"time"

	"api"
	"control/okex"
	mo "models/okex"
	"server/wshb"
	"util/config"
	"util/log"
	"util/sign"

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
	handlers *api.Handlers
	closeSig chan bool
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

func (app *App) CotrolHandlers() {

	app.closeSig = make(chan bool, 1)

	klGate := wshb.NewGate("wss://real.okex.com:10440/websocket/okexapi", 1, 1024, 65536, 5*time.Second, 5*time.Second, true, okex.NewAgentLogin(256))

	go klGate.Run(app.closeSig)

	time.Sleep(2000)

	var lg mo.ReqFurtureLogin
	timeStamp := fmt.Sprintf("%v", float32(time.Now().UnixNano()/1000))
	fmt.Println(timeStamp)
	lg.Event = "login"
	lg.Params.ApiKey = "342d1884-db81-4a9c-8535-1d4351965adf"
	lg.Params.PassPhrase = "IMDANDAN"
	lg.Params.Sign = sign.HMacSha256(timeStamp+"GET"+"/users/self/verify", []byte("3628818392EC421EF456070057E0F9CF"))
	lg.Params.TimeStamp = timeStamp

	klGate.Agent.WriteMsg(lg)
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
