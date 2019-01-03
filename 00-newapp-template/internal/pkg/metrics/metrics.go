package metrics

import (
	"00-newapp-template/pkg/acme"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type Metrics struct {
	server  serverMetric
	client  clientMetric
	service serviceMetric
}

type serverMetric struct {
	endPoint *prometheus.CounterVec
	cache    *prometheus.CounterVec
	db       *prometheus.CounterVec
}

type clientMetric struct {
	GophersThings prometheus.CounterVec
}

type serviceMetric struct {
	Transport *prometheus.CounterVec
}

func NewMetrics() (m *Metrics) {
	m = new(Metrics)

	m.serverInit()
	m.serviceInit()

	return m
}

func (m *Metrics) serverInit() {
	m.server.endPoint = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: `gophercli_server_endpoint_total`,
			Help: "How many service calls",
		}, []string{"method", "endpoint"},
	)
	m.server.cache = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: `gophercli_server_cache_total`,
			Help: "How many cache calls",
		}, []string{"method", "endpoint"},
	)
	m.server.db = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: `gophercli_server_db_total`,
			Help: "How many DB calls",
		}, []string{"method", "endpoint"},
	)
}

func (m *Metrics) serviceInit() {
	m.service.Transport = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: `gophercli_service_transport_total`,
			Help: "How many services calls",
		}, []string{"method", "endpoint"},
	)
}

func (m *Metrics) Marshal(filename string) {
	prometheus.WriteToTextfile(filename, prometheus.DefaultGatherer)
}

func (m *Metrics) ServerInc(endPoint acme.EndPoint, method serviceMethodType) {
	labels := prometheus.Labels{"endpoint": endPoint.String(), "method": method.String()}
	m.server.endPoint.With(labels).Inc()
}

func (m *Metrics) DBInc(endPoint acme.EndPoint, method dbMethodType) {
	labels := prometheus.Labels{"endpoint": endPoint.String(), "method": method.String()}
	m.server.db.With(labels).Inc()
}

func (m *Metrics) CacheInc(endPoint acme.EndPoint, method cacheMethodType) {
	labels := prometheus.Labels{"endpoint": endPoint.String(), "method": method.String()}
	m.server.cache.With(labels).Inc()
}

func (m *Metrics) TransportInc(endPoint acme.EndPoint, method transportMethodType) {
	labels := prometheus.Labels{"endpoint": endPoint.String(), "method": method.String()}
	m.service.Transport.With(labels).Inc()
}
