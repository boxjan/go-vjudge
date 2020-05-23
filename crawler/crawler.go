package crawler

import (
	"boxjan.li/go-vjudge/crawler/httpclient"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"sync"
	"time"
)

type Crawler interface {
	SupportOjList() []string
	RemoteOjSupportLanguage(remoteOjName string) ([]string, error)
	GetProblem(remoteOjName, problemId string) (*Problem, error)
	SubmitProblem(remoteOJName, problemId, code, language string) (*SolutionStatus, error)
	SolutionStatusUpdate(solutionStatus *SolutionStatus) error
}

type Problem struct {
	Title string
	TimeLimit uint
	MemoryLimit uint
	Description string
	Input string
	Output string
	Sample string
	Hint string
	Source string
	OriUrl string
}

type SolutionStatus struct {
	OriOj string
	OriProblemId string
	OriSolutionId string
	SubmitUser string
	UpdateAt time.Time
	Status RemoteStatus
	RawStatus string
	StatusHelper string
	MemoryUsed uint64
	TimeUsed uint64
}

type crawler struct {
	configPath string
	Config configStruct
	CreateAt time.Time
	workerHandle map[string]RemoteOj
	submitLock map[string]*sync.Mutex
	clients *httpclient.ClientsPool
	sync.Once
}

type Info interface {
	BaseUrl() string
	LanguageFinder() *map[string]string
}

type Submitter interface {
	GetMaxRunningId(*httpclient.Client, string, string) string
	Submit(*httpclient.Client, string, string, string) error
}

type ProblemCrawler interface {
	GetProblemUrl(string) string
	GetProblemContext(string) (*Problem, error)
}

type Login interface {
	Login(*httpclient.Client, *remoteAuth) error
	IsLogin(*httpclient.Client) bool
}

type StatusCrawler interface {
	GetRemoteStatus(*httpclient.Client, *SolutionStatus) error
}

type RemoteOj interface {
	Info
	Submitter
	ProblemCrawler
	Login
	StatusCrawler
}

func NewCrawler(path string) (*crawler, error) {
	handle := crawler{CreateAt:time.Now(),submitLock: map[string]*sync.Mutex{}}
	if 	s, err := os.Stat(path); err != nil || s.IsDir() {
		return nil, ErrConfigFileNotFound
	}
	handle.configPath = path
	handle.loadOjList()
	handle.reload()
	return &handle, nil
}

func (h *crawler)loadOjList() {
	if h.workerHandle == nil {
		h.workerHandle = make(map[string]RemoteOj)
	}

	{
		h.workerHandle["poj"] = NewPoj()
	}
}

func (h *crawler)SupportOjList() []string {
	res := make([]string, 0, len(h.workerHandle))
	for i := range h.workerHandle {
		res = append(res, i)
	}
	return res
}

func (h *crawler)RemoteOjSupportLanguage(OjName string) ([]string, error) {
	handle, err := h.getHandle(OjName)
	if err != nil {
		return nil, err
	}

	res := make([]string, 0, len(*(handle.LanguageFinder())))
	for lang := range *handle.LanguageFinder() {
		res = append(res, lang)
	}
	return res, nil
}

func (h *crawler)GetProblem(OjName, ProblemId string) (*Problem, error) {
	return h.getProblem(OjName, ProblemId)
}

func (h *crawler)SubmitProblem(remoteOjName string, problemId, code, language string ) (*SolutionStatus, error) {
	handle, err := h.getHandle(remoteOjName)
	if err != nil {
		return nil, err
	}
	auth := h.getAuthInfo(remoteOjName)
	cli, err := h.getLoginClient(remoteOjName, auth)
	if err != nil {
		return nil, err
	}

	if h.submitLock[remoteOjName + auth.Identity] == nil {
		h.submitLock[remoteOjName + auth.Identity] = &sync.Mutex{}
	}

	h.submitLock[remoteOjName + auth.Identity].Lock()
	defer h.submitLock[remoteOjName + auth.Identity].Unlock()

	beforeSubmitRunningId := handle.GetMaxRunningId(cli, auth.Identity, problemId)
	err = handle.Submit(cli, problemId, code, language)
	if err != nil {
		return nil, err
	}

	afterSubmittedRunningId := handle.GetMaxRunningId(cli, auth.Identity, problemId)

	if beforeSubmitRunningId == afterSubmittedRunningId {
		return nil, ErrSubmitFail
	}

	res := &SolutionStatus{
		OriOj:         remoteOjName,
		OriProblemId:  problemId,
		OriSolutionId: afterSubmittedRunningId,
		SubmitUser:    auth.Identity,
		UpdateAt:      time.Now(),
		Status:        StatusQueuing,
		StatusHelper:  "",
	}

	return res, nil
}

func (h *crawler)SolutionStatusUpdate(ss *SolutionStatus) error {
	handle, err := h.getHandle(ss.OriOj)
	if err != nil {
		return err
	}

	cli, err := h.getClientByOjNameAndIdentity(ss.OriOj, ss.SubmitUser)
	if err != nil {
		return err
	}
	err = handle.GetRemoteStatus(cli, ss)
	if err != nil {
		return err
	}
	return nil
}


func (h *crawler)getLoginClient(remoteOjName string, auth *remoteAuth) (cli *httpclient.Client, err error) {
	var handle RemoteOj
	handle, err = h.getHandle(remoteOjName)
	if err != nil {
		return
	}
	cli, err = h.getClientByOjNameAndIdentity(remoteOjName, auth.Identity)
	if err != nil {
		return
	}
	err = handle.Login(cli, auth)
	return
}


func (h * crawler)getHandle(ojName string) (RemoteOj, error ){
	crawlerHandle, ok := h.workerHandle[ojName]
	if !ok {
		return nil, ErrNotSupportOj
	}
	return crawlerHandle, nil
}

func (h *crawler)getProxyInfoForOj(ojName string) string {
	lenOfOjProxy := len(h.Config.RemoteMapByName[ojName].Proxies)
	if  lenOfOjProxy != 0 {
		return h.Config.RemoteMapByName[ojName].Proxies[(rand.Int() % lenOfOjProxy)]
	} else {
		return h.Config.Proxies[(rand.Int() % len(h.Config.Proxies))]
	}
}

func (h *crawler)getAuthInfo(ojName string) *remoteAuth {
	return &h.Config.RemoteMapByName[ojName].Accounts[(rand.Int() %len(h.Config.RemoteMapByName[ojName].Accounts))]
}

func (h *crawler)getProblem(remoteOjName, ProblemId string)  (problem *Problem, err error)  {
	var crawlerHandle RemoteOj
	crawlerHandle, err = h.getHandle(remoteOjName)
	if err != nil {
		return
	}

	problemUrl := crawlerHandle.GetProblemUrl(ProblemId)

	var rsp *httpclient.SimpleResponse
	var cli *httpclient.Client
	cli, err = h.getClientByOjName(remoteOjName)
	if err != nil {
		return
	}

	rsp = httpclient.TurnSimpleResponse(cli.Get(problemUrl))
	if rsp.Err != nil {
		err = rsp.Err
		return
	}

	if rsp.StatusCode < 200 || rsp.StatusCode > 299 {
		return nil, errors.New(http.StatusText(rsp.StatusCode))
	}

	problem, err =  crawlerHandle.GetProblemContext(rsp.Html)
	if err != nil {
		return
	}
	problem.OriUrl = problemUrl
	return
}

func (h* crawler)getClientByOjName(remoteOjName string) (*httpclient.Client, error) {
	return h.clients.GetClientByName(remoteOjName + "-defaultConn", h.getProxyInfoForOj(remoteOjName))
}

func (h *crawler)getClientByOjNameAndIdentity(remoteOjName, identity string) (*httpclient.Client, error) {
	return h.clients.GetClientByName(fmt.Sprintf("%s-%s", remoteOjName, identity), h.getProxyInfoForOj(remoteOjName))
}

func (h *crawler)NewSimpleRequestWithDefault(method, url string, postFrom, header map[string]string) *httpclient.SimpleRequest {
	if header == nil {
		header = make(map[string]string)
	}
	if _, ok := header["User-Agent"]; !ok {
		header["User-Agent"] = h.Config.Ua
	}
	return &httpclient.SimpleRequest{
		Method:   method,
		Url:      url,
		PostFrom: postFrom,
		Header:   header,
	}
}
func parseImages(urlToGet *url.URL, content string) ([]string, error) {
	var (
		err        error
		imgs       []string
		matches    [][]string
		findImages = regexp.MustCompile("<img.*?src=\"(.*?)\"")
	)

	// Retrieve all image URLs from string
	matches = findImages.FindAllStringSubmatch(content, -1)

	for _, val := range matches {
		var imgUrl *url.URL

		// Parse the image URL
		if imgUrl, err = url.Parse(val[1]); err != nil {
			return imgs, err
		}

		// If the URL is absolute, add it to the slice
		// If the URL is relative, build an absolute URL
		if imgUrl.IsAbs() {
			imgs = append(imgs, imgUrl.String())
		} else {
			imgs = append(imgs, urlToGet.Scheme+"://"+urlToGet.Host+imgUrl.String())
		}
	}

	return imgs, err
}