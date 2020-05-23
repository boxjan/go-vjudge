// Author: Boxjan
// Datetime: 2020/3/21 14:38
// I don't know why, but it can work

package crawler

import "errors"

var ErrConfigFileNotFound = errors.New("can not get config file in given path")
var ErrNotSupportOj = errors.New("not support oj name" )
var ErrRemoteOjReport = errors.New("remote OJ report a error")
var ErrSubmitFail = errors.New("submit failed")