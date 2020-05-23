package controllers

import (
	"boxjan.li/go-vjudge/models"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"html/template"
	"net/http"
)


type BaseController struct {
	beego.Controller
	User *models.User
	error map[string]string
	old map[string]string
}

func (c *BaseController) Prepare()  {
	c.Data["xsrf_token"] = c.XSRFToken()
	c.Data["csrf"]=template.HTML(c.XSRFFormHTML())

	var ok bool
	if c.old, ok = c.GetSession("old").(map[string]string); ok {
		c.Data["old"] = c.old
	} else {
		c.Data["old"] = make(map[string]string)
	}

	if c.error, ok = c.GetSession("error").(map[string]string); ok {
		c.Data["error"] = c.error
	} else {
		c.Data["error"] = make(map[string]string)
	}

	c.error = make(map[string]string)
	c.old = make(map[string]string)

	userIdRaw  := c.GetSession("user_id")
	if userId, ok := userIdRaw.(uint64); ok {
		user := models.User{Id: userId}
		err := orm.NewOrm().Read(&user)
		if err == nil {
			c.User = &user
		}
	}

	c.Data["user"] = c.User
	c.Layout = ""
	//c.LayoutSections = make(map[string]string)
}

func (c *BaseController)NeedLogin() {
	if c.User == nil {
		c.SetSession("url", c.Ctx.Request.URL.String())
		c.Redirect(beego.URLFor("LoginController.LoginPage"), http.StatusSeeOther)
		c.StopRun()
	}
}

func (c *BaseController)NoLogin() {
	if c.User != nil {
		c.Redirect(beego.URLFor("IndexController.Index"), http.StatusSeeOther)
		c.StopRun()
	}
}

func (c *BaseController)SetOld(key, value string) {
	c.old[key] = value
}

func (c *BaseController)SetError(key, err string) {
	c.error[key] = err
}

func (c *BaseController)GetOld(key string) string {
	if value, ok := c.old[key]; ok {
		return value
	}
	return ""
}

func (c *BaseController)GetErr(key string) string {
	if err, ok := c.error[key]; ok {
		return err
	}
	return ""
}

func (c *BaseController)Finish()  {
	c.SetSession("error", c.error)
	c.SetSession("old", c.old)
}