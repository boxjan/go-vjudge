// Author: Boxjan
// Datetime: 2020/3/19 16:55
// I don't know why, but it can work

package crawler

import (
	"boxjan.li/go-vjudge/crawler/httpclient"
	"boxjan.li/go-vjudge/tools"
	"encoding/base64"
	"errors"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

type poj struct {}

var pojLanguageFinder = map[string]string{
	"G++": "0",
	"GCC": "1",
	"Java": "2",
	"Pascal": "3",
	"C++": "4",
	"C": "5",
	"Fortran": "6",
	}

func NewPoj() *poj {
	return &poj{}
}

func (oh *poj)BaseUrl() string {
	return "http://poj.org/"
}

func (oh *poj)LanguageFinder() *map[string]string { // showed => submit
	return &pojLanguageFinder
}

func (oh* poj)GetProblemUrl(problemId string) string {
	return oh.BaseUrl() + "problem?id=" + problemId
}

func (oh *poj)GetProblemContext(html string) (*Problem, error) {
	timeLimitStr := tools.RegexFindOneWithoutError(html, "<b>Time Limit:</b> (\\d{3,})MS</td>")
	timeLimit, err := strconv.Atoi(timeLimitStr)
	if err != nil {
		return nil, err
	}
	memLimitStr := tools.RegexFindOneWithoutError(html, "<b>Memory Limit:</b> (\\d{3,})K</td>")
	memLimit, err := strconv.Atoi(memLimitStr)
	if err != nil {
		return nil, err
	}
	return &Problem{
		Title:       tools.RegexFindOneWithoutError(html, "<title>\\d{3,} -- ([\\s\\S]*?)</title>"),
		TimeLimit:   uint(timeLimit),
		MemoryLimit: uint(memLimit),
		Description: tools.RegexFindOneWithoutError(html, "<p class=\"pst\">Description</p>([\\s\\S]*?)<p class=\"pst\">"),
		Input:       tools.RegexFindOneWithoutError(html, "<p class=\"pst\">Input</p>([\\s\\S]*?)<p class=\"pst\">"),
		Output:      tools.RegexFindOneWithoutError(html, "<p class=\"pst\">Output</p>([\\s\\S]*?)<p class=\"pst\">"),
		Sample:      tools.RegexFindOneWithoutError(html, "(<p class=\"pst\">Sample Input</p>[\\s\\S]*?Sample Output[\\s\\S]*?)<p class=\"pst\">"),
		Hint:        tools.RegexFindOneWithoutError(html, "<p class=\"pst\">Hint</p>([\\s\\S]*?)<p class=\"pst\">"),
		Source:      tools.RegexFindOneWithoutError(html, "<p class=\"pst\">Source</p>([\\s\\S]*?)</td></tr></tbody></table>"),
	}, nil
}

func (oh *poj)IsLogin(c *httpclient.Client) bool {
	rsp := httpclient.TurnSimpleResponse(c.Get(oh.BaseUrl()))
	if rsp.Err != nil {
		return false
	}
	return strings.Contains(rsp.Html, ">Log Out</a>")
}

func (oh *poj)Login(client *httpclient.Client, auth *remoteAuth) error {
	if oh.IsLogin(client) {
		return nil
	}

	rsp, err := client.PostForm(oh.BaseUrl() + "login", map[string][]string{
		"B1": {"login"},
		"password1": {auth.Password},
		"url": {"/"},
		"user_id1": {auth.Identity},
	})
	if err != nil {
		return err
	}

	if rsp.StatusCode != http.StatusFound {
		return errors.New("login failed")
	}

	if !oh.IsLogin(client) {
		return errors.New("login meet problem")
	}

	return nil
}

func (oh *poj)GetMaxRunningId(cl *httpclient.Client, identify, problemId string) string {
	rsp := httpclient.TurnSimpleResponse(cl.Get(oh.BaseUrl() + "status?user_id=" + identify + "&problem_id=" + problemId))
	if rsp.Err != nil {
		return ""
	}
	return tools.RegexFindOneWithoutError(rsp.Html, "<tr align=center><td>(\\d+)")

}

func (oh *poj)Submit(cl *httpclient.Client, problemId, code, language string) error {

	rsp := httpclient.TurnSimpleResponse(cl.PostForm(oh.BaseUrl() + "submit", map[string][]string{
		"language": {(*oh.LanguageFinder())[language]},
		"problem_id": {problemId},
		"source": {base64.StdEncoding.EncodeToString([]byte(code))},
		"encoded": {"1"},
	}))
	if rsp.Err != nil {
		return rsp.Err
	}
	if rsp.StatusCode != http.StatusFound {
		return errors.New("submit failed")
	}
	return nil
}

func (oh *poj)GetRemoteStatus(cli *httpclient.Client, status *SolutionStatus) error {
	rsp := httpclient.TurnSimpleResponse(cli.Get(oh.BaseUrl() +  "/showsource?solution_id=" + status.OriSolutionId))
	if rsp.Err != nil {
		return rsp.Err
	}
	if strings.Contains(rsp.Html, "<title>Error</title>") {
		return ErrRemoteOjReport
	}

	var err error
	rawStatus := tools.RegexFindOneWithoutError(rsp.Html, "<b>Result:</b>(.+?)</td>")
	status.RawStatus = regexp.MustCompile("<.*?>").ReplaceAllString(rawStatus, "")
	status.Status = RemoteStatusType(status.RawStatus)
	if  status.Status == StatusCompileError {
		if rsp := httpclient.TurnSimpleResponse(cli.Get(oh.BaseUrl() + "/showcompileinfo?solution_id=" + status.OriSolutionId)); rsp.Err != nil {
			return rsp.Err
		} else {
			status.StatusHelper = tools.RegexFindOneWithoutError(rsp.Html,"<pre>([\\s\\S]*?)</pre>")
		}
	} else if status.Status == StatusAccept {
		status.MemoryUsed, err = strconv.ParseUint(tools.RegexFindOneWithoutError(rsp.Html,  "<b>Memory:</b> ([-\\d]+)"), 10, 10)
		if err != nil {
			return err
		}
		status.TimeUsed, err = strconv.ParseUint(tools.RegexFindOneWithoutError(rsp.Html,  "<b>Time:</b> ([-\\d]+)"), 10, 10)
		if err != nil {
			return err
		}
	}
	return nil
}