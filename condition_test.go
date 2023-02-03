package cond

import (
	"strconv"
	"testing"
)

type Event struct {
	Type  string
	Value int
}

func (ev *Event) Field(key string) string {
	switch key {
	case "type":
		return ev.Type
	case "value":
		return strconv.Itoa(ev.Value)
	}

	return ""
}

func TestExp(t *testing.T) {
	cnd := New("value ~ (.*)")
	ev := &Event{
		Type:  "typeof",
		Value: 456,
	}

	pay := func(id int, ret string) {
		t.Logf("%d %v", id, ret)
	}

	t.Log(cnd.Match(ev, Payload(pay)))
}

func TestString(t *testing.T) {

	raw := "12-345-67.raw"

	pbc := String(raw)
	ext := pbc("[:6]")

	t.Log(ext)

}
func TestRegex(t *testing.T) {
	val := "10.10.239.11"
	cnd := New("[0,13] ~ \\.(.*)\\.(.*)\\.(.*)")

	pay := func(id int, ret string) {
		t.Logf("%d %v", id, ret)
	}

	t.Log(cnd.Match(val, Payload(pay)))
}
