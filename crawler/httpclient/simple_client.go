package httpclient

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type Client struct {
	Header http.Header
	Client *http.Client
}

type SimpleRequest struct {
	Method string
	Url string
	PostFrom map[string]string
	Header map[string]string
}

type SimpleResponse struct {
	StatusCode int
	Html string
	Err error
}

func DoSimpleRequest(client *Client, request *SimpleRequest) *SimpleResponse {
	if !strings.HasPrefix(request.Url, "http://")  && !strings.HasPrefix(request.Url, "https://") {
		request.Url = "http://" + request.Url
	}

	PostFrom := make(url.Values)
	for k, v := range request.PostFrom {
		PostFrom[k] = []string{v}
	}

	req, err := http.NewRequest(request.Method, request.Url, strings.NewReader(PostFrom.Encode()))
	if err != nil {
		return &SimpleResponse{Err:err}
	}

	for k, v := range request.Header {
		req.Header.Set(k, v)
	}
	if request.Method == "POST" {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}
	return TurnSimpleResponse(client.Do(req))
}

func TurnSimpleResponse(rsp *http.Response, sourceErr error) *SimpleResponse  {
	if sourceErr != nil {
		return &SimpleResponse{Err:sourceErr}
	}
	responseHtmlByte, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		log.Println(err)
	}
	_ = rsp.Body.Close()
	return &SimpleResponse{
		StatusCode: rsp.StatusCode,
		Html:       string(responseHtmlByte),
	}
}


func (sp *SimpleResponse)String() string {
	if sp.Err != nil {
		return sp.Err.Error()
	}
	return "Status: " + http.StatusText(sp.StatusCode) + " Html:" + sp.Html
}

func (c *Client) Get(url string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	for k, v := range c.Header {
		if len(req.Header.Get(k)) == 0 {
			req.Header[k] = append(req.Header[k], v...)
		}
	}
	return c.Client.Do(req)
}

func (c *Client) Head(url string) (*http.Response, error) {
	req, err := http.NewRequest("HEAD", url, nil)
	if err != nil {
		return nil, err
	}
	for k, v := range c.Header {
		if len(req.Header.Get(k)) == 0 {
			req.Header[k] = append(req.Header[k], v...)
		}
	}
	return c.Client.Do(req)
}


func (c *Client) Post(url, contentType string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", contentType)
	for k, v := range c.Header {
		if len(req.Header.Get(k)) == 0 {
			req.Header[k] = append(req.Header[k], v...)
		}
	}
	return c.Client.Do(req)
}

func (c *Client) PostForm(url string, data url.Values) (*http.Response, error) {
	return c.Post(url, "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))
}

func (c *Client) Do(req *http.Request) (*http.Response, error) {
	for k, v := range c.Header {
		if len(req.Header.Get(k)) == 0 {
			req.Header[k] = append(req.Header[k], v...)
		}
	}
	return c.Client.Do(req)
}

func (c *Client) DoSimpleRequest(req *SimpleRequest) *SimpleResponse {
	return DoSimpleRequest(c, req)
}