package middleware

import (
	"net/http"
	"ratelimit/infra/redis"
	"time"
)

func RateLimiterMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		key := r.URL.Path + "_127.0.0.1"
		if !RateLimiter.ShouldAllow(key) {
			http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})

}

// NOTE: This does not require manual expiry and is preferable since data is persistent in redis.
// This allows to have throttleCount(no of request) in throttleCount(time limit). If there are more requests error is thrown.
func GetThrottlingMiddleWare(throttleDuration time.Duration, throttleCount int64) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			redisKey := r.URL.Path + "_127.0.0.1"

			callCount := redis.Incr(redisKey)

			if callCount >= throttleCount {
				http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
			}

			if callCount <= 1 {
				redis.SetExpiry(&ctx, redisKey, throttleDuration)
			}

			next.ServeHTTP(w, r)
		})
	}
}
