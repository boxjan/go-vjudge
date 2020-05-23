package models

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"net/url"
)
var o orm.Ormer

func init()  {
	var dsm, t string
	switch beego.AppConfig.String("dbType") {
	case "pgsql": fallthrough
	case "postgres": {
		url, _ := url.Parse("postgres://" + beego.AppConfig.String("dbAddr"))
		dsm = fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
			url.Hostname(), url.Port(), beego.AppConfig.String("dbUser"),
			beego.AppConfig.String("dbPass"), beego.AppConfig.String("dbName"))
		t = "postgres"
	}

	case "mysql": {
		dsm = fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			beego.AppConfig.String("dbUser"), beego.AppConfig.String("dbPass"),
			beego.AppConfig.String("dbAddr"), beego.AppConfig.String("dbName"))
		t = "mysql"
	}


	case "sqlite":{
		dsm = fmt.Sprintf("%s", beego.AppConfig.String("dbName"))
		t = "sqlite"
	}
	default:
		panic("not support database type: " + beego.AppConfig.String("dbType"))
	}

	orm.RegisterDataBase("default", t, dsm, 30)

	// register model
	orm.RegisterModel(new(User), new(Problem), new(ProblemDescription),
		new(Solution), new(SolutionCode), new(SolutionStatusHelper))
	o = orm.NewOrm()
}

