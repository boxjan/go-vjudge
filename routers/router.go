package routers

import (
	"boxjan.li/go-vjudge/controllers"
	"boxjan.li/go-vjudge/controllers/auth"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
)


var FilterMethod = func(ctx *context.Context) {
	if ctx.Input.Query("_method")!="" && ctx.Input.IsPost(){
		ctx.Request.Method = ctx.Input.Query("_method")
	}
}

func init() {
	beego.InsertFilter("*", beego.BeforeRouter, FilterMethod)

    beego.Router("/", &controllers.IndexController{}, "*:Index")

	beego.Router("/login", &auth.LoginController{},"get:LoginPage;post:Login")
	beego.Router("/register", &auth.RegisterController{},"get:RegisterPage;post:Register")
	beego.Router("/logout", &auth.LoginController{}, "get:Logout;post:Logout")

	//beego.Router("/user/info", &controllers.UserController{}, "get:InfoPage;post:UpdateInfo")
	//beego.Router("/user/resetpassword", &controllers.UserController{}, "get:ResetPasswordPage;post:ResetPassword")

	beego.Router("/problem", &controllers.ProblemController{}, "get:List")
	beego.Router("/problem/:id", &controllers.ProblemController{}, "get:Detail")
	beego.Router("/api/problem", &controllers.ProblemController{}, "get:ApiList")
	beego.Router("/api/problem/:id", &controllers.ProblemController{}, "get:ApiDetail")


	beego.Router("/submit/", &controllers.SolutionController{}, "Get:UseInRoute")
	beego.Router("/submit/:id", &controllers.SolutionController{}, "Get:SubmitPage;Post:Submit")
	beego.Router("/solution", &controllers.SolutionController{}, "get:List")
	beego.Router("/api/solution", &controllers.SolutionController{}, "get:ApiList")
	beego.Router("/api/solution/:id/error", &controllers.SolutionController{}, "get:ApiErrorDetail")
	beego.Router("/api/solution/:id/code", &controllers.SolutionController{}, "get:ApiCode")


}
