package controllers

import (
	"boxjan.li/go-vjudge/crawler"
	"boxjan.li/go-vjudge/models"
	"github.com/astaxie/beego/orm"
	"log"
	"strings"
	"time"
)

type ProblemStruct struct {
	models.Problem
	models.ProblemDescription
}

type ProblemIntro struct {
	Id uint64 `json:"id"`
	OriOj string `json:"ori_oj"`
	OriId string `json:"ori_id"`
	Title string `json:"title"`
	Source string `json:"source"`
	UpdateAt time.Time `json:"update_at"`
}

type ProblemDataForApi struct {
	Search *ProblemSearch `json:"search"`
	Cur int `json:"cur"`
	Tol int64 `json:"tol"`
	List []ProblemIntro `json:"list"`
}

type ProblemSearch struct {
	OriOj string `json:"ori_oj" form:"ori_oj"`
	OriId string `json:"ori_id" form:"ori_id"`
	Title string `json:"title" form:"title"`
	Source string `json:"source" form:"search"`
	Page int `json:"page" form:"page"`
}

type ProblemController struct {
	BaseController
}

func (c *ProblemController)List()  {
	var problemSearch ProblemSearch
	err := c.ParseForm(&problemSearch)
	if err != nil {
		c.Abort("404")
	}

	c.Data["supportOJ"] = crawler.SupportOjList()
	c.Data["search"] = problemSearch
	c.TplName="problem/list.tpl"
}

func (c *ProblemController)ApiList()  {
	var problems []ProblemIntro
	var problemSearch ProblemSearch
	err := c.ParseForm(&problemSearch)
	if err != nil {
		c.Abort("404")
	}
	if strings.ToLower(problemSearch.OriOj) == "all" {
		problemSearch.OriOj = ""
	}
	if problemSearch.Page == 0 {
		problemSearch.Page = 1
	}

	var problemsModel []*models.Problem
	var count, allCount int64
	count, err = orm.NewOrm().QueryTable(models.ProblemTableName).
		Filter("OriOj__contains", problemSearch.OriOj).
		Filter("OriProblemId__startswith", problemSearch.OriId).
		Filter("Title__contains", problemSearch.Title).
		Filter("Source__contains", problemSearch.Source).
		OrderBy("-UpdatedAt").Limit(15).Offset((problemSearch.Page - 1) * 15).All(&problemsModel)
	if err != nil {
		log.Println(err)
		c.Abort("500")
	}
	allCount, err = orm.NewOrm().QueryTable(models.ProblemTableName).
		Filter("OriOj__contains", problemSearch.OriOj).
		Filter("OriProblemId__startswith", problemSearch.OriId).
		Filter("Title__contains", problemSearch.Title).
		Filter("Source__contains", problemSearch.Source).
		OrderBy("-UpdatedAt").Count()
	if err != nil {
		log.Println(err)
		c.Abort("500")
	}

	if count == 0 {
		if problemSearch.OriOj != "" && problemSearch.OriId != "" {
			crawler.GetProblemThenSaveToDb(problemSearch.OriOj, problemSearch.OriId)
		}
	} else {
		problems = make([]ProblemIntro, count)
		for i := int64 (0); i < count; i++ {
			problems[i].Id = problemsModel[i].Id
			problems[i].Title = problemsModel[i].Title
			problems[i].Source = problemsModel[i].Source
			problems[i].UpdateAt = problemsModel[i].UpdatedAt
			problems[i].OriOj = problemsModel[i].OriOj
			problems[i].OriId = problemsModel[i].OriProblemId
		}
	}


	c.Data["json"] = &ProblemDataForApi{
		Search: &problemSearch,
		Cur:  problemSearch.Page,
		Tol:  allCount,
		List: problems,
	}
	c.ServeJSON()
}

func (c *ProblemController)Detail()  {
	problemId, err:= c.GetUint64(":id", 0)
	if err != nil || problemId == 0 {
		c.Abort("404")
	}

	problem := ProblemStruct{}
	problem.Id = problemId
	err = orm.NewOrm().Read(&problem.Problem)
	if err != nil {
		log.Print(err)
		c.Abort("404")
	}
	var problemDescription *models.ProblemDescription
	problemDescription, err = problem.Problem.GetProblemDesc()
	if err != nil {
		log.Print(err)
		c.Abort("404")
	}
	problem.ProblemDescription = *problemDescription

	c.Data["problem"] = &problem
	c.Data["problem_id"] = problemId
	c.TplName="problem/detail.tpl"
}

func (c *ProblemController)ApiDetail()  {
	problemId, err:= c.GetUint64(":id", 0)
	if err != nil || problemId == 0 {
		c.Abort("404")
	}

	var problem ProblemStruct


	c.Data["json"] = &problem
	c.ServeJSON()
}

