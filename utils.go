package hiperus

import (
	"github.com/ziutek/soap"
)

func asBool(e *soap.Element) bool {
	s := e.AsStr()
	return s == "t" || s == "true"
}
