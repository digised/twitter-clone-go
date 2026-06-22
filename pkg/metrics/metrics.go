package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var HTTPRequestsTotal = promauto.NewCounterVec(
	prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Total number of HTTP requests",
	},
	[]string{"method", "path", "status"},
)

var HTTPRequestDuration = promauto.NewHistogramVec(
	prometheus.HistogramOpts{
		Name:    "http_request_duration_seconds",
		Help:    "HTTP request duration in seconds",
		Buckets: prometheus.DefBuckets, // 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10
	},
	[]string{"method", "path"},
)

var UsersCreatedTotal = promauto.NewCounter(prometheus.CounterOpts{
	Name: "users_created_total",
	Help: "Total number of registered users",
})

var FollowsTotal = promauto.NewCounter(prometheus.CounterOpts{
	Name: "follows_total",
	Help: "Total number of follow actions",
})

var UnfollowsTotal = promauto.NewCounter(prometheus.CounterOpts{
	Name: "unfollows_total",
	Help: "Total number of unfollow actions",
})

var FollowRateLimitHitsTotal = promauto.NewCounter(prometheus.CounterOpts{
	Name: "follow_rate_limit_hits_total",
	Help: "Number of times the follow rate limit was exceeded",
})
