package middleware

import (
	"encoding/json"
	"context"
	"net/http"
	"time"
	"fmt"
	"os"

	"github.com/mileusna/useragent"
)

type Information struct {
	Timestamp time.Time `json:"timestamp"`	
	Latency int64 `json:"latency"`
	Path string `json:"path"`
	OS string `json:"os"`
}

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

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		next.ServeHTTP(w,r)

		end := time.Now()

		latency := end.Sub(start).Milliseconds()

		path := r.URL.Path

		os, ok := r.Context().Value(OSKey).(string)
		if !ok {
			os = "unknown"
		}

		logData := Information{
			Timestamp: start,
			Latency: latency,
			Path: path,
			OS: os,
		}

		jsonData, err := json.Marshal(logData)
		if err != nil {
			fmt.Println("json error:", err)
			return
		}

		fmt.Println(string(jsonData))
	})
}

func BasicAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userId, password, ok := r.BasicAuth()

		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		expectedUserId := os.Getenv("BASIC_AUTH_USER_ID")
		expectedPassword := os.Getenv("BASIC_AUTH_PASSWORD")

		if userId != expectedUserId || password != expectedPassword {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w,r)

	})
}