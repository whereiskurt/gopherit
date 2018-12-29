package pkg

import (
	"00-newapp-template/pkg/acme"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type Metrics struct {
	Config MetricsConfig

	Server  ServerMetric
	Client  ClientMetric
	Service ServiceMetric
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

	m.InitServerMetric()
	m.InitServiceMetric()

	return m
}

func (m *Metrics) InitServerMetric() {
	m.Server.EndPoint = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: `gophercli_server_endpoint_total`,
			Help: "How many service calls",
		}, []string{"method", "endpoint"},
	)
	m.Server.Cache = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: `gophercli_server_cache_total`,
			Help: "How many cache calls",
		}, []string{"method", "endpoint"},
	)
	m.Server.DB = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: `gophercli_server_db_total`,
			Help: "How many DB calls",
		}, []string{"method", "endpoint"},
	)
}

func (m *Metrics) InitServiceMetric() {
	m.Service.Transport = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: `gophercli_service_transport_total`,
			Help: "How many services calls",
		}, []string{"method", "endpoint"},
	)
}

func (m *Metrics) EndPointInc(endPoint acme.ServiceEndPoint, method string) {
	labels := prometheus.Labels{"endpoint": fmt.Sprintf("%s", endPoint), "method": method}
	m.Server.EndPoint.With(labels).Inc()
}
func (m *Metrics)  DBInc(endPoint acme.ServiceEndPoint, method string) {
	labels := prometheus.Labels{"endpoint": fmt.Sprintf("%s", endPoint), "method": method}
	m.Server.DB.With(labels).Inc()
}
func (m *Metrics)  CacheInc(endPoint acme.ServiceEndPoint, method string) {
	labels := prometheus.Labels{"endpoint": fmt.Sprintf("%s", endPoint), "method": method}
	m.Server.Cache.With(labels).Inc()
}