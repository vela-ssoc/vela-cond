package cond

import (
	"bytes"
	"github.com/vela-ssoc/vela-kit/lua"
	"github.com/vela-ssoc/vela-kit/pipe"
	"github.com/vela-ssoc/vela-kit/vela"
	"net"
	"strings"
)

var xEnv vela.Environment

/*
	cnd = vela.cnd.compile("name eq zhangsan,lisi,wangwu").ok(a).no(b)
    vela.cnd.like("abc")


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

func newLuaCidrJit(L *lua.LState) int {
	n := L.GetTop()
	if n < 2 {
		L.RaiseError("invalid cidr condition args")
		return 0
	}

	cnd := &LCond{
		onMatch: pipe.New(pipe.Env(xEnv)),
		noMatch: pipe.New(pipe.Env(xEnv)),
		co:      xEnv.Clone(L),
	}

	keys := strings.Split(L.CheckString(1), ",")
	s := &Section{
		keys:   keys,
		method: Cidr,
	}

	for i := 2; i <= n; i++ {
		_, ipNet, err := net.ParseCIDR(L.CheckString(i))
		if err != nil {
			L.RaiseError("invalid parse ip net #%d", i)
			return 0
		}
		s.subnet = append(s.subnet, ipNet)
	}

	cnd.cnd = &Cond{
		data: []*Section{s},
	}
	L.Push(cnd)
	return 1
}

func WithEnv(env vela.Environment) {
	xEnv = env

	kv := lua.NewUserKV()
	kv.Set("format", lua.NewFunction(newLuaConditionFormat))
	kv.Set("cidr", lua.NewFunction(newLuaCidrJit))
	kv.Set("group", lua.NewFunction(NewCombineL))

	xEnv.Set("cnd", lua.NewExport("vela.condition.export", lua.WithTable(kv), lua.WithFunc(newLuaConditionJit)))
	xEnv.Set("cndf", lua.NewFunction(newLuaConditionFormat))
}

/*

vela.cnd.cidr("src,dst" , 192.168.0.0/24","10.0.0.0/24")
vela.cnd("abc").match(object)




*/
