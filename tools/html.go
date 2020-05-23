// Author: Boxjan
// Datetime: 2020/3/28 21:14
// I don't know why, but it can work

package tools

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

const DefaultCharset = "utf-8"

func GetCharset(rsp *http.Response) string {
	fromHeaderRaw := rsp.Header.Get("Content-Type")
	fromHeader := strings.ToLower(RegexFindOneWithoutError(fromHeaderRaw, "charset=([\\s\\S]*);?"))

	if len(fromHeader) != 0 {
		return fromHeader
	}

	responseHtmlByte, err := ioutil.ReadAll(rsp.Body)
	if err == nil {
		fromHtml := strings.ToLower(RegexFindOneWithoutError(string(responseHtmlByte), "<meta charset=\"(.*)?\">"))
		if len(fromHtml) != 0 {
			return fromHtml
		}
	}
	return DefaultCharset
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