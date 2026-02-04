package middleware

import "net/http"

type PanicHandler struct{}

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