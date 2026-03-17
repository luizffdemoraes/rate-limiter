package middleware

import (
	"net"
	"net/http"
	"strings"

	"github.com/lffm1994/rate-limiter/internal/limiter"
)

const tooManyRequestsBody = "you have reached the maximum number of requests or actions allowed within a certain time frame"

func RateLimiter(l *limiter.Limiter) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := clientIP(r)
			token := strings.TrimSpace(r.Header.Get("API_KEY"))

			res, err := l.Allow(ip, token)
			if err != nil {
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}

			if !res.Allowed {
				w.WriteHeader(http.StatusTooManyRequests)
				_, _ = w.Write([]byte(tooManyRequestsBody))
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func clientIP(r *http.Request) string {
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		parts := strings.Split(xff, ",")
		return strings.TrimSpace(parts[0])
	}

	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}
