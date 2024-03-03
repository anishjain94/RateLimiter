package middleware

import (
	"fmt"
	"sync"
	"time"
)

type UserRequestCount struct {
	Lock *sync.Mutex

	LastVisited   time.Time
	RequestsCount uint64
}

type RateLimitConfig struct {
	MaxRequests uint64
	Expiry      time.Duration
}

type RateLimit struct {
	Lock *sync.Mutex

	Config          map[string]*RateLimitConfig  //Here key would be APIPATH_IP
	Users           map[string]*UserRequestCount //Here key would be APIPATH_IP
	CleanupInterval time.Duration
}

var RateLimiter RateLimit

func (rl RateLimit) cleanup() {
	for range time.Tick(rl.CleanupInterval) {
		for key, user := range rl.Users {
			config := rl.Config[key]
			if config != nil && time.Since(user.LastVisited) > config.Expiry {
				fmt.Println("deleted:", key)
				rl.Lock.Lock()
				delete(rl.Users, key)
				rl.Lock.Unlock()
			}
		}
	}
}

func AddConfig(key string, rate uint64, expiry time.Duration) {

	RateLimiter.Lock.Lock()
	defer RateLimiter.Lock.Unlock()

	RateLimiter.Config[key] = &RateLimitConfig{
		MaxRequests: rate,
		Expiry:      expiry,
	}

}

func NewRateLimiter(cleanup time.Duration) {
	RateLimiter = RateLimit{
		Config:          make(map[string]*RateLimitConfig),
		Users:           make(map[string]*UserRequestCount),
		CleanupInterval: cleanup,
		Lock:            &sync.Mutex{},
	}

	go RateLimiter.cleanup()

}

func (rateLimit RateLimit) ShouldAllow(key string) bool {

	rateLimit.Lock.Lock()
	defer rateLimit.Lock.Unlock()

	userReuqestCount, exists := rateLimit.Users[key]
	if !exists {
		rateLimit.Users[key] = &UserRequestCount{
			LastVisited:   time.Now(),
			RequestsCount: 1,
			Lock:          &sync.Mutex{},
		}
		return true
	}

	config, exists := rateLimit.Config[key]
	if config == nil {
		return true
	}

	if !exists || time.Since(userReuqestCount.LastVisited) > config.Expiry {
		rateLimit.Users[key] = &UserRequestCount{
			Lock:          &sync.Mutex{},
			LastVisited:   time.Now(),
			RequestsCount: 1,
		}

		return true
	}

	userReuqestCount.Lock.Lock()
	defer userReuqestCount.Lock.Unlock()

	if userReuqestCount.RequestsCount < config.MaxRequests {
		userReuqestCount.RequestsCount++
		userReuqestCount.LastVisited = time.Now()
		return true
	}

	return false

}
