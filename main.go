package main

import (
	"boxjan.li/go-vjudge/crawler"
	_ "boxjan.li/go-vjudge/models"
	_ "boxjan.li/go-vjudge/routers"
	"boxjan.li/go-vjudge/tools"
	"boxjan.li/go-vjudge/views"
	"encoding/base64"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

func main() {
	views.InitRenderFunc()
	if beego.AppConfig.String("runmode") == "dev" {
		orm.Debug = true
		orm.RunSyncdb("default", false, true)
	}
	if len(beego.AppConfig.String("xsrfkey")) == 0 {
		key := base64.StdEncoding.EncodeToString(tools.RandBytes(32))
		beego.AppConfig.Set("xsrfkey", key)
		beego.AppConfig.SaveConfigFile("conf/app.conf")
	}

	beego.SetStaticPath("/ojFiles", "static/ojFiles")

	crawler.InitCrawler(beego.AppConfig.String("crawlerconfig"))

	beego.Run()
}

