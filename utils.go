package hiperus

import (
	"errors"
	"github.com/ziutek/soap"
)

func asBool(e *soap.Element) bool {
	s := e.AsStr()
	return s == "t" || s == "true"
}

var ErrEmptyResultSet = errors.New("hiperus: empty result set")

func firstRow(rs *soap.Element) (*soap.Element, error) {
	if len(rs.Children) == 0 {
		return nil, ErrEmptyResultSet
	}
	return rs.Children[0], nil
}
