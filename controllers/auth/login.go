package auth

import (
	"boxjan.li/go-vjudge/controllers"
	"boxjan.li/go-vjudge/models"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

type LoginForm struct {
	Identification string `form:"identification"`
	Password string `form:"password"`
}

type LoginController struct {
	controllers.BaseController
}

func (c *LoginController)LoginPage()  {
	c.NoLogin()
	if url, ok := c.GetSession("url").(string); ok {
		if len(url) > 0 {
			c.Data["url"] = true
		} else {
			c.Data["url"] = false
		}
	}
	c.Data["title"] = "login"
	c.TplName = "auth/login.tpl"
}

func (c *LoginController)Login()   {
	c.NoLogin()

	loginInfo := LoginForm{}
	if err := c.ParseForm(&loginInfo); err != nil {
		c.SetError("identification","Username or password wrong")
		c.Redirect(beego.URLFor("LoginController.LoginPage"), http.StatusSeeOther)
		return
	}

	user := models.User{Username: loginInfo.Identification}
	c.SetOld("identification", loginInfo.Identification)

	if err := orm.NewOrm().Read(&user, "Username"); err != nil {
		c.SetError("identification","Username or password wrong")
		c.Redirect(beego.URLFor("LoginController.LoginPage"), http.StatusSeeOther)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginInfo.Password)); err != nil {
		c.SetError("identification","Username or password wrong")
		c.Redirect(beego.URLFor("LoginController.LoginPage"), http.StatusSeeOther)
		return
	}

	c.SetSession("user_id", user.Id)

	if url, ok :=c.GetSession("url").(string); ok {
		c.DelSession("url")
		c.Redirect(url, http.StatusSeeOther)
		return
	}
	c.Redirect(beego.URLFor("IndexController.Index"), http.StatusSeeOther)
}

func (c *LoginController)Logout()   {
	c.NeedLogin()

	c.User = nil
	c.DelSession("user_id")
	c.DestroySession()
	c.Redirect(beego.URLFor("IndexController.Index"), http.StatusSeeOther)
}
