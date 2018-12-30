package pkg

import (
	"00-newapp-template/pkg/acme"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type Metrics struct {
	Config MetricsConfig

	server  ServerMetric
	client  ClientMetric
	service ServiceMetric
}

type ServerMetric struct {
	EndPoint *prometheus.CounterVec
	Cache    *prometheus.CounterVec
	DB       *prometheus.CounterVec
}

type ClientMetric struct {
	GophersThings prometheus.CounterVec
}

type ServiceMetric struct {
	Transport *prometheus.CounterVec
}

func NewMetrics(config MetricsConfig) (m *Metrics) {
	m = new(Metrics)

	m.initServerMetric()
	m.initServiceMetric()

	return m
}

func (m *Metrics) initServerMetric() {
	m.server.EndPoint = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: `gophercli_server_endpoint_total`,
			Help: "How many service calls",
		}, []string{"method", "endpoint"},
	)
	m.server.Cache = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: `gophercli_server_cache_total`,
			Help: "How many cache calls",
		}, []string{"method", "endpoint"},
	)
	m.server.DB = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: `gophercli_server_db_total`,
			Help: "How many DB calls",
		}, []string{"method", "endpoint"},
	)
}

func (m *Metrics) initServiceMetric() {
	m.service.Transport = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: `gophercli_service_transport_total`,
			Help: "How many services calls",
		}, []string{"method", "endpoint"},
	)
}

func (m *Metrics) TransportInc(endPoint acme.ServiceEndPoint, method string) {
	labels := prometheus.Labels{"endpoint": fmt.Sprintf("%s", endPoint), "method": method}
	m.service.Transport.With(labels).Inc()
}
func (m *Metrics) EndPointInc(endPoint acme.ServiceEndPoint, method string) {
	labels := prometheus.Labels{"endpoint": fmt.Sprintf("%s", endPoint), "method": method}
	m.server.EndPoint.With(labels).Inc()
}
func (m *Metrics) CacheInc(endPoint acme.ServiceEndPoint, method string) {
	labels := prometheus.Labels{"endpoint": fmt.Sprintf("%s", endPoint), "method": method}
	m.server.Cache.With(labels).Inc()
}
func (m *Metrics) DBInc(endPoint acme.ServiceEndPoint, method string) {
	labels := prometheus.Labels{"endpoint": fmt.Sprintf("%s", endPoint), "method": method}
	m.server.DB.With(labels).Inc()
}
