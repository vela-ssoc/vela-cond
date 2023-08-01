package cond

import (
	"net/url"
)

const (
	Eq op = iota + 10
	Re
	Cn
	In
	Lt
	Le
	Ge
	Gt
	Unary
	Call
	Oop
	Pass
	Regex
	Cidr
)

var (
	opTab = []string{"equal", "grep", "contain", "include", "less", "less or equal", "greater or equal", "greater", "unary", "call", "oop", "pass", "regex"}
)

type op uint8

func (o op) call(v string, raw string) bool {
	if len(v) == 0 {
		return false
	}

	uri, err := url.Parse(raw)
	if err != nil {
		return false
	}

	c := &TnlCall{
		uri: uri,
	}
	return c.do(v)
}

func (o op) String() string {
	return opTab[(int(o) - 10)]
}
