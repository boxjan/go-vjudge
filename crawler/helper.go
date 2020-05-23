package crawler

import (
	"boxjan.li/go-vjudge/models"
	timeLRU "boxjan.li/go-vjudge/tools/time_lru"
	"github.com/astaxie/beego/orm"
	"log"
	"sync"
	"time"
)

type Worker struct {
	sync.Once
	handle Crawler
	lru *timeLRU.TimeLRU
}

var worker Worker

func InitCrawler(config string)  {
	worker.Do(func() {
		var err error
		worker.handle, err = NewCrawler(config)
		if err != nil {
			panic(err)
		}
		worker.lru = timeLRU.New(128)
	})
}

func SupportOjList() []string {
	return worker.handle.SupportOjList()
}

func RemoteOjSupportLanguage(remoteOjName string) ([]string, error) {
	return worker.handle.RemoteOjSupportLanguage(remoteOjName)
}

func GetProblem(remoteOjName, problemId string) (*Problem, error) {
	return worker.handle.GetProblem(remoteOjName, problemId)
}

func SubmitProblem(remoteOJName, problemId, code, language string) (*SolutionStatus, error) {
	return worker.handle.SubmitProblem(remoteOJName, problemId, code, language)
}

func SolutionStatusUpdate(solutionStatus *SolutionStatus) error {
	return worker.handle.SolutionStatusUpdate(solutionStatus)
}

func GetProblemThenSaveToDb(remoteOjName, remoteProblemId string)  {

	if _, ok := worker.lru.Get(remoteOjName + remoteProblemId); ok {
		return
	}
	problem, err := GetProblem(remoteOjName, remoteProblemId)

	if err != nil  {
		worker.lru.Add(remoteOjName + remoteProblemId, "not found", 10 * time.Second)
		log.Printf("get problem: %s %s meet error %+v", remoteOjName, remoteProblemId)
		return
	}
	//return
	problemModel := models.Problem{}
	//var problemId int64

	err = orm.NewOrm().QueryTable(models.ProblemTableName).Filter("OriOj", remoteOjName).Filter("OriProblemId", remoteProblemId).One(&problemModel)
	if err != nil && err != orm.ErrNoRows{
		log.Println(err)
		return
	}
	if err == orm.ErrNoRows {
		NewProblem := models.Problem{
			OriOj:        remoteOjName,
			OriProblemId: remoteProblemId,
			OriUrl:       problem.OriUrl,
			Title:        problem.Title,
			Source:       problem.Source,
			MemoryLimit:  problem.MemoryLimit,
			TimeLimit:    problem.TimeLimit,
		}
		problemId, err := orm.NewOrm().Insert(&NewProblem)
		if err != nil {
			log.Println(err)
			return
		}
		NewProblem.Id = uint64(problemId)
		err = NewProblem.InsertProblemDesc(&models.ProblemDescription{
			Description: problem.Description,
			Input:       problem.Input,
			Output:      problem.Output,
			Sample:      problem.Sample,
			Hint:        problem.Hint,
		})
		if err != nil {
			log.Println(err)
			return
		}
	} else {
		problemModel.UpdatedAt = time.Now()
		orm.NewOrm().Update(&problemModel, "updated_at")
		err = problemModel.UpdateProblemDesc(&models.ProblemDescription{
			Description: problem.Description,
			Input:       problem.Input,
			Output:      problem.Output,
			Sample:      problem.Sample,
			Hint:        problem.Hint,
		})
		if err != nil {
			log.Println(err)
			return
		}
	}

}

func SubmitProblemFormDb(solutionId uint64) {
	go func() {
		solution := models.Solution{Id:solutionId}
		err := orm.NewOrm().Read(&solution)
		if err != nil {
			log.Println(err)
			return
		}
		solutionCode, err := solution.GetCode()
		if err != nil {
			log.Println(err)
			return
		}
		solutionStatus, err :=SubmitProblem(solution.OriOj, solution.OriProblemId, solutionCode.Code, solution.Language)
		if err != nil {
			log.Println(err)
			return
		}
		solution.OriSubmitTime = solutionStatus.UpdateAt
		solution.OriAccount = solutionStatus.SubmitUser
		solution.OriSolutionId = solutionStatus.OriSolutionId
		_, _ = orm.NewOrm().Update(&solution)
		for solutionStatus.Status < 0 {
			err = SolutionStatusUpdate(solutionStatus)
			if err != nil {
				log.Println(err)
				return
			}
			solution.Status = int (solutionStatus.Status)
			solution.RawStatus = solutionStatus.RawStatus
			_, err = orm.NewOrm().Update(&solution, "Status", "RawStatus")
			if err != nil {
				log.Println(err)
				return
			}
			time.Sleep(500 * time.Millisecond)
		}

		solution.Status = int (solutionStatus.Status)
		solution.RawStatus = solutionStatus.RawStatus
		solution.OriAccount = solutionStatus.SubmitUser
		solution.TimeUsed = uint(solutionStatus.TimeUsed)
		solution.MemoryUsed = uint(solutionStatus.MemoryUsed)

		_, err = orm.NewOrm().Update(&solution)
		if err != nil {
			log.Println(err)
			return
		}
		err = solution.InsertStatusHelper(&models.SolutionStatusHelper{
			StatusHelper: solutionStatus.StatusHelper,
		})
		if err != nil {
			log.Println(err)
			return
		}
	}()
}