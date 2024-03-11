package cond

import (
	"bytes"
	"github.com/vela-ssoc/vela-kit/auxlib"
	"github.com/vela-ssoc/vela-kit/lua"
)

const (
	OR Logic = iota + 1
	AND
)

type Logic uint8

func (l Logic) String() string {
	switch l {
	case OR:
		return "or"
	case AND:
		return "and"
	}
	return "unknown"
}

type Cond struct {
	data []*Section
}

func New(c ...string) *Cond {
	n := len(c)
	if n == 0 {
		return &Cond{}
	}

	cond := &Cond{
		data: make([]*Section, len(c)),
	}

	for i := 0; i < n; i++ {
		cond.data[i] = Compile(c[i])
	}
	return cond
}

func F(prefix string, v ...interface{}) string {
	if len(v) == 0 {
		return ""
	}

	var buf bytes.Buffer
	buf.Write(auxlib.S2B(prefix))
	buf.WriteByte(' ')
	for i, item := range v {
		if i != 0 {
			buf.WriteByte(',')
		}
		buf.WriteString(auxlib.ToString(item))
	}

	return buf.String()

}

func CheckMany(L *lua.LState, opt ...OptionFunc) *Cond {
	cnd := &Cond{}
	cnd.CheckMany(L, opt...)
	return cnd
}

func Check(L *lua.LState, idx int) *LCond {
	ov := L.CheckObject(idx)

	lc, ok := ov.(*LCond)
	if ok {
		lc.co = xEnv.Clone(L)
		return lc
	}

	L.RaiseError("invalid condition object , got %p", &ov)
	return nil
}

func LValue(L *lua.LState, val lua.LValue) *LCond {
	if val.Type() != lua.LTObject {
		L.RaiseError("invalid condition type , got %v", val.Type().String())
	}

	lc, ok := val.(*LCond)
	if ok {
		lc.co = xEnv.Clone(L)
		return lc
	}

	L.RaiseError("invalid condition object , got %p", &val)
	return nil
}
