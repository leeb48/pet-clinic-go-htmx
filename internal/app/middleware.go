package app

import "net/http"

func (app *App) LogRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			ip     = r.RemoteAddr
			proto  = r.Proto
			method = r.Method
			uri    = r.URL.RequestURI()
			body   = r.Body
		)

		app.Logger.Info("request", "ip", ip, "proto", proto, "method", method, "uri", uri, "body", body)
		next.ServeHTTP(w, r)
	})
}
