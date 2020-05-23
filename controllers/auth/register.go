package auth

import (
	"boxjan.li/go-vjudge/controllers"
	"boxjan.li/go-vjudge/models"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/validation"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strings"
)

type RegisterInfo struct {
	Username string `form:"username" valid:"Required; MaxSize(60)"`
	Nickname string `form:"nickname" valid:"Required; MaxSize(60)"`
	Email string `form:"email" valid:"Required; Email; MaxSize(250)"`
	Password string `form:"password" valid:"Required"`
	PasswordConfirm string `form:"password_confirmation" valid:"Required"`
	School string `form:"school" valid:"Required; MaxSize(120)"`
}

type RegisterController struct {
	controllers.BaseController
}

func (c *RegisterController)RegisterPage()  {
	c.NoLogin()
	c.TplName = "auth/register.tpl"
}

func (c *RegisterController)Register()  {
	c.NoLogin()

	registerInfo := RegisterInfo{}
	if err := c.ParseForm(&registerInfo); err != nil {
		logs.Warn("register meet error: ", err)
		c.Redirect(beego.URLFor("RegisterController.RegisterPage"), http.StatusSeeOther)
		return
	}

	c.SetOld("username", registerInfo.Username)
	c.SetOld("nickname", registerInfo.Nickname)
	c.SetOld("email", registerInfo.Email)
	c.SetOld("school", registerInfo.School)

	valid := validation.Validation{}
	if b, err := valid.Valid(&registerInfo); err != nil {
		logs.Warn("valid meet error: ", err)
		c.Redirect(beego.URLFor("RegisterController.RegisterPage"), http.StatusSeeOther)
		return
	} else if !b {
		for _, err := range valid.Errors {
			key := strings.ToLower(strings.Split(err.Key, ".")[0])
			c.SetError(key, err.Error())
		}
		c.Redirect(beego.URLFor("RegisterController.RegisterPage"), http.StatusSeeOther)
		return
	}

	if registerInfo.Password != registerInfo.PasswordConfirm {
		c.SetError("password", "Password is not same asPassword Confirm")
		c.Redirect(beego.URLFor("RegisterController.RegisterPage"), http.StatusSeeOther)
		return
	}

	if count, err := orm.NewOrm().QueryTable(models.UserTableName).Filter("username", registerInfo.Username).Count(); err != nil {
		logs.Warn("sql meet error: ", err)
		c.SetError("username", "something error, please try again")
		c.Redirect(beego.URLFor("RegisterController.RegisterPage"), http.StatusSeeOther)
		return
	} else if count > 0 {
		c.SetError("username", "the username have been used")
		c.Redirect(beego.URLFor("RegisterController.RegisterPage"), http.StatusSeeOther)
		return
	}

	passwordWithSalt, err := bcrypt.GenerateFromPassword([]byte(registerInfo.Password), 5)
	if err != nil {
		logs.Warn("password add salt meet error: ", err)
		c.SetError("password", "something error, please try again")
		c.Redirect(beego.URLFor("RegisterController.RegisterPage"), http.StatusSeeOther)
	}
	user := models.User{
		Username:        registerInfo.Username,
		Nickname:        registerInfo.Nickname,
		Password:        string(passwordWithSalt),
		Role:            "user",
		Email:           registerInfo.Email,
		School:          registerInfo.School,
		RegisterIp:      c.Ctx.Input.IP(),
	}

	if id, err := orm.NewOrm().Insert(&user); err != nil {
		logs.Warn("sql meet error: ", err)
		c.SetError("username", "something error, please try again")
		c.Redirect(beego.URLFor("RegisterController.RegisterPage"), http.StatusSeeOther)
	} else {
		c.SetSession("user_id", id)
	}
	c.Redirect(beego.URLFor("IndexController.Index"), http.StatusSeeOther)
}