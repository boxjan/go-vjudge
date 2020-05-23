package controllers

import (
	"boxjan.li/go-vjudge/crawler"
	"boxjan.li/go-vjudge/models"
	"crypto/sha256"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/validation"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"
)

type SolutionController struct {
	BaseController
}

type SolutionSubmit struct {
	ProblemId uint64 `json:"problem_id" form:"problem_id" valid:"Required"`
	Code string `json:"code" form:"code" valid:"Required"`
	Language string `json:"language" form:"language" valid:"Required"`
}

type SolutionSearch struct {
	OriOj string `json:"ori_oj" form:"ori_oj"`
	OriId string `json:"ori_id" form:"ori_id"`
	User string `json:"user" form:"user"`
	Status string `json:"status" form:"status"`
	Page int `json:"page" form:"page"`
}

type SolutionIntro struct {
	SolutionId uint64 `json:"solution_id"`
	User string `json:"user"`
	OriOj string `json:"ori_oj"`
	OriId string `json:"ori_id"`
	ProblemId uint64 `json:"problem_id"`
	Status int `json:"status"`
	RawStatus string `json:"raw_status"`
	MemoryUsed uint `json:"memory_used"`
	TimeUsed uint `json:"time_used"`
	Lang string `json:"lang"`
	Length uint `json:"length"`
	SubmitTime time.Time `json:"submit_time"`
}

type SolutionDataForApi struct {
	Search *SolutionSearch `json:"search"`
	Cur int `json:"cur"`
	Tol int64 `json:"tol"`
	List []SolutionIntro `json:"list"`
}

func (c* SolutionController)UseInRoute()  {
	c.Abort("404")
}

func (c* SolutionController)SubmitPage()  {
	c.NeedLogin()
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
	langs, err := crawler.RemoteOjSupportLanguage(problem.OriOj)
	if err != nil {
		log.Print(err)
		c.Abort("404")
	}
	sort.Strings(langs)

	lastLang, ok := c.GetSession(problem.OriOj + "last_lang").(string)
	if !ok {
		lastLang = ""
	}

	c.Data["problem_id"] = problemId
	c.Data["problem_title"] = problem.Title
	c.Data["problem_ori_oj"] = problem.OriOj
	c.Data["language"] = langs
	c.Data["last_lang"] = lastLang
	c.TplName="solution/submit.tpl"
}

func (c* SolutionController)Submit()  {
	c.NeedLogin()
	problemId, err := c.GetUint64(":id", 0)
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

	submit := SolutionSubmit{}
	err = c.ParseForm(&submit)
	if err != nil {
		c.Abort("404")
	}
	c.SetSession(problem.OriOj + "last_lang", submit.Language)
	c.SetOld("code", submit.Code)

	valid := validation.Validation{}
	if b, err := valid.Valid(&submit); err != nil {
		logs.Warn("valid meet error: ", err)
		c.Redirect(beego.URLFor("SolutionController.UseInRoute") + "/" + strconv.FormatUint(problemId,10), http.StatusSeeOther)
		return
	} else if !b {
		for _, err := range valid.Errors {
			key := strings.ToLower(strings.Split(err.Key, ".")[0])
			c.SetError(key, err.Error())
		}
		c.Redirect(beego.URLFor("SolutionController.UseInRoute") + "/" + strconv.FormatUint(problemId,10), http.StatusSeeOther)
		return
	}


	hash := fmt.Sprintf("%x", sha256.Sum256([]byte(strings.TrimSpace(submit.Code))))


	solution := models.Solution{
		UserId:        c.User.Id,
		ProblemId:     submit.ProblemId,
		Language:      submit.Language,
		Length: 	   uint(len(submit.Code)),
		OriOj:         problem.OriOj,
		OriProblemId:  problem.OriProblemId,
		Hash:          string(hash[:]),
		RawStatus:     "",
		Status:        int (crawler.StatusPending),
	}


	solution_id, err := orm.NewOrm().Insert(&solution)
	if err != nil {
		log.Println(err)
		c.SetError("code", "we just meet some error, please try again")
		c.Redirect(beego.URLFor("SolutionController.UseInRoute") + "/" + strconv.FormatUint(problemId,10), http.StatusSeeOther)
		return
	}
	solution.Id = uint64(solution_id)
	err = solution.InsertCode(&models.SolutionCode{Code:submit.Code})
	if err != nil {
		log.Println(err)
		return
	}

	crawler.SubmitProblemFormDb(solution.Id)
	c.Redirect(beego.URLFor("SolutionController.List") +
		"?user=" + c.User.Username,
	302)
}

func (c *SolutionController)List()  {
	var solutionSearch SolutionSearch
	err := c.ParseForm(&solutionSearch)
	if err != nil {
		c.Abort("404")
	}

	c.Data["supportOJ"] = crawler.SupportOjList()
	c.Data["search"] = &solutionSearch
	c.TplName="solution/list.tpl"
}

func (c *SolutionController)ApiList()  {
	var solutions []SolutionIntro
	var solutionSearch SolutionSearch
	err := c.ParseForm(&solutionSearch)
	if err != nil {
		c.Abort("404")
	}
	if strings.ToLower(solutionSearch.OriOj) == "all" {
		solutionSearch.OriOj = ""
	}
	if solutionSearch.Page == 0 {
		solutionSearch.Page = 1
	}
	var solutionModel []*models.Solution
	var count, allCount int64

	queryS := orm.NewOrm().QueryTable(models.SolutionTableName).
		Filter("OriOj__contains", solutionSearch.OriOj).
		Filter("OriProblemId__startswith", solutionSearch.OriId)

	if solutionSearch.Status != "" {
		if i, err := strconv.Atoi(solutionSearch.Status); err == nil {
			queryS = queryS.Filter("Status", i)
		}
	}
	if solutionSearch.User != "" {
		userId := username2UserId(solutionSearch.User)
		if userId != 0 {
			queryS = queryS.Filter("UserId", userId)
		} else {
			c.Abort("404")
		}
	}

	count, err = queryS.OrderBy("-UpdatedAt").Limit(15).Offset((solutionSearch.Page - 1) * 15).All(&solutionModel)
	if err != nil {
		log.Println(err)
		c.Abort("500")
	}
	allCount, err = queryS.Count()
	if err != nil {
		log.Println(err)
		c.Abort("500")
	}

	solutions = make([]SolutionIntro, 0, count)
	for  _, solution := range solutionModel {
		solutions = append(solutions, SolutionIntro{
			SolutionId: solution.Id,
			User:       userId2username(solution.UserId),
			OriOj:      solution.OriOj,
			OriId:      solution.OriProblemId,
			ProblemId:  solution.ProblemId,
			Status:     solution.Status,
			RawStatus:  solution.RawStatus,
			MemoryUsed: solution.MemoryUsed,
			TimeUsed:   solution.TimeUsed,
			Lang:       solution.Language,
			Length:     solution.Length,
			SubmitTime: solution.CreatedAt,
		})
	}
	c.Data["json"] = &SolutionDataForApi{
		Search: &solutionSearch,
		Cur:  solutionSearch.Page,
		Tol:  allCount,
		List: solutions,
	}
	c.ServeJSON()
}

func (c *SolutionController)ApiErrorDetail()  {
	solutionId, err := c.GetUint64(":id", 0)
	if err != nil || solutionId == 0 {
		c.CustomAbort(404, "Not Found")
	}
	ssh := models.SolutionStatusHelper{SolutionId:solutionId}
	err = orm.NewOrm().Read(&ssh, "solution_id")
	if err != nil {
		c.CustomAbort(404, "Not Found")
	}
	c.CustomAbort(200, ssh.StatusHelper)
}

func (c *SolutionController)ApiCode() {
	c.NeedLogin()
	solutionId, err := c.GetUint64(":id", 0)
	if err != nil || solutionId == 0 {
		c.CustomAbort(404, "Not Found")
	}
	solution := models.Solution{Id:solutionId}
	err = orm.NewOrm().Read(&solution)
	if err != nil {
		c.CustomAbort(404, "Not Found")
	}
	if solution.UserId != c.User.Id {
		c.CustomAbort(404, "Not Found")
	}
	sc , err := solution.GetCode()
	if err != nil {
		c.CustomAbort(404, "Not Found")
		return
	}
	c.CustomAbort(200, sc.Code)
}

func username2UserId(username string) uint64 {
	user := models.User{Username:username}
	err := orm.NewOrm().Read(&user, "username")
	if err == nil {
		return user.Id
	}
	return 0
}

func userId2username(userId uint64) string {
	user := models.User{Id:userId}
	err := orm.NewOrm().Read(&user)
	if err == nil {
		return user.Username
	}
	return ""
}