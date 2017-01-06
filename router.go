package goweb

import "regexp"

type router struct {
	re      *regexp.Regexp
	handler HTTPHandler
}
