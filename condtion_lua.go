package cond

import (
	"github.com/vela-ssoc/vela-kit/lua"
	"github.com/vela-ssoc/vela-kit/pipe"
)

type LCond struct {
	onMatch *pipe.Px
	noMatch *pipe.Px
	co      *lua.LState
	cnd     *Cond
}

func (lc *LCond) String() string                         { return "vela.condition" }
func (lc *LCond) Type() lua.LValueType                   { return lua.LTObject }
func (lc *LCond) AssertFloat64() (float64, bool)         { return 0, false }
func (lc *LCond) AssertString() (string, bool)           { return "", false }
func (lc *LCond) AssertFunction() (*lua.LFunction, bool) { return nil, false }
func (lc *LCond) Peek() lua.LValue                       { return lc }

func (lc *LCond) okL(L *lua.LState) int {
	lc.onMatch = pipe.NewByLua(L, pipe.Seek(0), pipe.Env(xEnv))
	L.Push(lc)
	return 1
}

func (lc *LCond) noL(L *lua.LState) int {
	lc.noMatch = pipe.NewByLua(L, pipe.Seek(0), pipe.Env(xEnv))
	L.Push(lc)
	return 1
}

func (lc *LCond) Match(lv lua.LValue, L *lua.LState) bool {
	if lc.cnd.Match(lv) {
		lc.onMatch.Do(lv, L, func(err error) {
			xEnv.Errorf("condition match function pipe fail %v", err)
		})

		return true
	}

	lc.noMatch.Do(lv, L, func(err error) {
		xEnv.Errorf("condition not match function pipe fail %v", err)
	})
	return false
}

func (lc *LCond) matchL(L *lua.LState) int {
	n := L.GetTop()
	if n == 0 {
		L.Push(lua.LFalse)
		return 1
	}

	for i := 1; i <= n; i++ {
		if lc.Match(L.Get(i), L) {
			L.Push(lua.LTrue)
			return 1
		}
	}

	L.Push(lua.LFalse)
	return 1
}

func (lc *LCond) Index(L *lua.LState, key string) lua.LValue {
	switch key {
	case "ok":
		return lua.NewFunction(lc.okL)
	case "no":
		return lua.NewFunction(lc.noL)
	case "match":
		return lua.NewFunction(lc.matchL)
	}

	return lua.LNil
}

func (lc *LCond) NewIndex(L *lua.LState, key string, val lua.LValue) {
	switch key {
	case "ok":
		lc.onMatch = pipe.New(pipe.Env(xEnv))
		lc.onMatch.LValue(val)
	case "no":
		lc.onMatch = pipe.New(pipe.Env(xEnv))
		lc.onMatch.LValue(val)
	}

}

func newLCond(L *lua.LState) *LCond {
	return &LCond{
		onMatch: pipe.New(pipe.Env(xEnv)),
		noMatch: pipe.New(pipe.Env(xEnv)),
		co:      xEnv.Clone(L),
		cnd:     CheckMany(L, Seek(0)),
	}
}
