package cond

import (
	"bytes"
	"github.com/vela-ssoc/vela-kit/lua"
	"github.com/vela-ssoc/vela-kit/vela"
)

var xEnv vela.Environment

/*
	cnd = vela.cnd.compile("name eq zhangsan,lisi,wangwu").ok(a).no(b)

	hash := "lisi"
	cnd.match("lisi") @end(true)

	vela.cnd.v("key = " , "ssss")
*/

func newLuaConditionJit(L *lua.LState) int {
	L.Push(newLCond(L))
	return 1
}

func newLuaConditionFormat(L *lua.LState) int {
	n := L.GetTop()
	if n < 2 {
		return 0
	}

	var buf bytes.Buffer
	prefix := L.IsString(1)
	buf.WriteString(prefix)
	buf.WriteByte(' ')

	for i := 2; i <= n; i++ {
		if i == 2 {
			buf.WriteByte(' ')
		} else {
			buf.WriteByte(',')
		}

		buf.WriteString(L.Get(i).String())
	}

	L.Push(lua.B2L(buf.Bytes()))
	return 1
}

func WithEnv(env vela.Environment) {
	xEnv = env
	xEnv.Set("cnd", lua.NewFunction(newLuaConditionJit))
	xEnv.Set("cndf", lua.NewFunction(newLuaConditionFormat))
}
