package middleware


import (
	"context"
	"net/http"

	"github.com/mileusna/useragent"

)

type contextKey string

const OSKey contextKey = "os"

type PanicHandler struct{}

func UserAgentMiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// User-Agentを取得
		uaString := r.UserAgent()

		// User-Agentを解析
		ua := useragent.Parse(uaString)
		os := ua.OS

		// User-Agent情報をコンテキストに保存
		ctx := context.WithValue(r.Context(), OSKey, os)

		// 新しいリクエストに差し替え
		r = r.WithContext(ctx)

		// 次のハンドラーを呼び出す
		next.ServeHTTP(w, r)


	})
}

func Recovery(h http.Handler) http.Handler {
fn := func(w http.ResponseWriter, r *http.Request) {
// TODO: ここに実装をする
defer func() {
	if rec := recover(); rec != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}()
h.ServeHTTP(w, r)
}
return http.HandlerFunc(fn)
}

func (ph PanicHandler) ServeHTTP(w http.ResponseWriter, r * http.Request) {
	panic("This is a panic example")
}