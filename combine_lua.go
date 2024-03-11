package cond

import (
	"github.com/vela-ssoc/vela-kit/kind"
	"github.com/vela-ssoc/vela-kit/lua"
	"github.com/vela-ssoc/vela-kit/strutil"
)

func (iv *Combine) String() string                         { return strutil.B2S(iv.Json()) }
func (iv *Combine) Type() lua.LValueType                   { return lua.LTObject }
func (iv *Combine) AssertFloat64() (float64, bool)         { return 0, false }
func (iv *Combine) AssertString() (string, bool)           { return "", false }
func (iv *Combine) AssertFunction() (*lua.LFunction, bool) { return lua.NewFunction(iv.Call), true }
func (iv *Combine) Peek() lua.LValue                       { return iv }

func (iv *Combine) Json() []byte {
	v := *iv
	n := len(v)
	if n == 0 {
		return []byte("[]")
	}

	enc := kind.NewJsonEncoder()
	enc.Arr("")
	for i := 0; i < n; i++ {
		enc.Val(v[i].String())
		enc.Val(",")
	}
	enc.End("]")
	return enc.Bytes()
}

func (iv *Combine) cndL(L *lua.LState) int {
	iv.CheckMany(L, WithCo(L))
	L.Push(iv)
	return 1
}

func (iv *Combine) Index(L *lua.LState, key string) lua.LValue {
	switch key {
	case "cnd":
		return lua.NewFunction(iv.cndL)
	default:
		return lua.LNil
	}
}

func (iv *Combine) Call(L *lua.LState) int {
	ret := iv.Match(L.Get(1), WithCo(L))
	L.Push(lua.LBool(ret))
	return 1
}

func NewCombineL(L *lua.LState) int {
	c := NewCombine()
	L.Push(c)
	return 1
}
