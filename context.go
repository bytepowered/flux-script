package script

import (
	"math/rand"
	"net/http"
	"net/url"
	"time"
)

import (
	"github.com/bytepowered/fluxgo/pkg/common"
	"github.com/bytepowered/fluxgo/pkg/flux"
	"github.com/bytepowered/fluxgo/pkg/logger"
	"github.com/spaolacci/murmur3"
)

type GetVarFunc func(key string) string

// ScriptContext 注入到JavaScript脚本引擎的上下文对象
type ScriptContext struct {
	RequestPattern string      `json:"pattern"`
	RequestMethod  string      `json:"method"`
	RequestPath    string      `json:"path"`
	RequestHost    string      `json:"host"`
	HeaderValues   http.Header `json:"headers"`
	FormValues     url.Values  `json:"forms"`
	QueryValues    url.Values  `json:"queries"`
	// Function
	GetPathVarFunc   GetVarFunc `json:"getPathVar"`
	GetQueryVarFunc  GetVarFunc `json:"getQueryVar"`
	GetHeaderVarFunc GetVarFunc `json:"getHeaderVar"`
	GetFormVarFunc   GetVarFunc `json:"getFormVar"`
	LookupExprFunc   GetVarFunc `json:"lookupExpr"`
	// Helper
	RandomInt63Func func(max int64) int64    `json:"random"`
	FastHashFunc    func(data string) uint64 `json:"hash"`
	ConsoleLogFunc  func(max interface{})    `json:"log"`
}

func NewScriptContext(webc flux.WebContext, pattern string) ScriptContext {
	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	return ScriptContext{
		RequestPattern: pattern,
		RequestMethod:  webc.Method(),
		RequestPath:    webc.URI(),
		RequestHost:    webc.Host(),
		HeaderValues:   webc.HeaderVars(),
		FormValues:     webc.FormVars(),
		QueryValues:    webc.QueryVars(),
		GetPathVarFunc: func(key string) string {
			return webc.PathVar(key)
		},
		GetQueryVarFunc: func(key string) string {
			return webc.QueryVar(key)
		},
		GetHeaderVarFunc: func(key string) string {
			return webc.HeaderVar(key)
		},
		GetFormVarFunc: func(key string) string {
			return webc.FormVar(key)
		},
		LookupExprFunc: func(expr string) string {
			return common.LookupWebValueByExpr(webc, expr)
		},
		RandomInt63Func: func(max int64) int64 {
			return random.Int63n(max)
		},
		FastHashFunc: func(data string) uint64 {
			return hash64([]byte(data))
		},
		ConsoleLogFunc: func(arg interface{}) {
			logger.Trace(webc.RequestId()).Info(arg)
		},
	}
}

func hash64(data []byte) uint64 {
	var h64 = murmur3.New64()
	_, _ = h64.Write(data)
	return h64.Sum64()
}
