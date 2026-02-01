package ratelimiter

import (
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/Avito-courses/course-go-avito-Grisha1Kadetov/iternal/metrics"
	"github.com/Avito-courses/course-go-avito-Grisha1Kadetov/iternal/pkg/log"
)

type bucket struct {
	tokens float64
	last   time.Time
}

type RateLimiter struct {
	Rps   float64
	Burst float64

	log     log.Logger
	mu      sync.Mutex
	buckets map[string]*bucket
}

func New(l log.Logger, rps, burst float64) *RateLimiter {
	return &RateLimiter{
		log:     l,
		buckets: make(map[string]*bucket),
		Rps:     rps,
		Burst:   burst,
	}
}

func (rl *RateLimiter) Middleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := getIP(r)

			if !rl.ipIsAllow(ip) {
				w.WriteHeader(http.StatusTooManyRequests)
				w.Header().Set("Retry-After", "1")
				rl.log.Warn("rate limit exceeded", log.NewField("ip", ip))

				metrics.RateLimitExceededTotal.WithLabelValues(
					r.URL.Path,
					r.Method,
				).Inc()

				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func (rl *RateLimiter) ipIsAllow(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	b, ok := rl.buckets[ip]
	if !ok {
		b = &bucket{
			tokens: rl.Burst,
			last:   time.Now(),
		}
		rl.buckets[ip] = b
	}

	now := time.Now()
	elapsed := now.Sub(b.last).Seconds()
	b.last = now

	b.tokens += elapsed * rl.Rps
	if b.tokens > rl.Burst {
		b.tokens = rl.Burst
	}

	if b.tokens < 1 {
		return false
	}

	b.tokens -= 1
	return true
}

func getIP(r *http.Request) string {
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		ip := strings.Split(xff, ",")[0]
		return strings.TrimSpace(ip)
	}

	if ip := r.Header.Get("X-Real-IP"); ip != "" {
		return ip
	}

	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err == nil {
		return host
	}

	return r.RemoteAddr
}
