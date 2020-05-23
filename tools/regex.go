// Author: Boxjan
// Datetime: 2020/3/28 19:54
// I don't know why, but it can work

package tools

import (
	"errors"
	"regexp"
)

var ErrRegexpNotFound = errors.New("not thing match")

func RegexFindOne(finding string, regex string) (string, error) {
	regexCompile := regexp.MustCompile(regex)
	params := regexCompile.FindSubmatch([]byte(finding))
	if len(params) <= 1 {
		return "", ErrRegexpNotFound
	}
		return string(params[1]), nil
}

func RegexFindOneWithoutError(finding string, regex string) string {
	res, _ := RegexFindOne(finding, regex)
	return res
}