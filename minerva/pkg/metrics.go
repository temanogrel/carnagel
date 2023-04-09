package minerva

import "github.com/prometheus/client_golang/prometheus"

var (
	grpcCallCounter = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "grpc",
		Name:      "call_counter",
	}, []string{"method"})

	downloadsRecommended = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "grpc",
		Name:      "recommended_downloads_counter",
	}, []string{"edge", "origin"})

	uploadsRecommended = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "grpc",
		Name:      "recommended_storage_counter",
	}, []string{"hostname"})
)

func init() {
	prometheus.Register(grpcCallCounter)
	prometheus.Register(downloadsRecommended)
	prometheus.Register(uploadsRecommended)
}

type metricReporter struct {
	GrpcCallCounter      *prometheus.CounterVec
	DownloadsRecommended *prometheus.CounterVec
	StorageRecommended   *prometheus.CounterVec
}

func GetMetricReporter() *metricReporter {
	return &metricReporter{
		GrpcCallCounter:      grpcCallCounter,
		StorageRecommended:   uploadsRecommended,
		DownloadsRecommended: downloadsRecommended,
	}
}
