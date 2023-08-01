package cond

import (
	"bytes"
	"fmt"
	"github.com/vela-ssoc/vela-kit/kind"
	"github.com/vela-ssoc/vela-kit/vela"
	"io"
	"net/http"
	"net/url"
	"strconv"
)

var (
	JsonHeader = http.Header{
		"Content-Type": []string{"application/json"},
	}
)

type TnlCall struct {
	uri   *url.URL
	reply struct {
		Count int                 `json:"count"`
		Data  map[string][]string `json:"data"`
	}
}

func (tc *TnlCall) Cache() vela.Bucket {
	cache := tc.uri.Query().Get("cache")
	if len(cache) == 0 {
		return nil
	}

	return xEnv.Bucket("cond", "cache", cache)
}

func (tc *TnlCall) Hit(db vela.Bucket, key string) bool {
	if db == nil {
		return false
	}

	v, err := db.Get(key)
	if err == nil {
		return false
	}

	if hit, ok := v.(bool); ok {
		return hit
	}

	return false
}

func (tc *TnlCall) Path() string {
	return fmt.Sprintf("/api/v1/broker/security%s?%s", tc.uri.Path, tc.uri.Query())
}

//func (tc *TnlCall) Query() string {
//	return tc.uri.RawQuery
//}

func (tc *TnlCall) Header() http.Header {
	return JsonHeader
}

func (tc *TnlCall) Body(raw string) io.Reader {
	enc := kind.NewJsonEncoder()
	enc.Tab("")
	enc.Join("data", []string{raw})
	enc.End("}")

	return bytes.NewReader(enc.Bytes())
}

func (tc *TnlCall) save(db vela.Bucket, key string, val bool) {
	if db == nil {
		return
	}

	ttl, _ := strconv.Atoi(tc.uri.Query().Get("ttl"))
	if ttl == 0 {
		ttl = 30
	}

	db.Store(key, val, ttl*1000)
}

func (tc *TnlCall) do(val string) bool {
	db := tc.Cache()

	if tc.Hit(db, val) {
		return true
	}

	if e := xEnv.JSON(tc.Path(), tc.Body(val), &tc.reply); e != nil {
		return false
	}

	if tc.reply.Count == 0 {
		tc.save(db, val, false)
		return false
	}

	tc.save(db, val, true)
	return true
}
