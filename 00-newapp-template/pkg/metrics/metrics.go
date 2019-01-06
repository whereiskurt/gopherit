package metrics

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type EndPointType string

var EndPoints = endPointTypes{
	Gophers: EndPointType("Gophers"),
	Gopher:  EndPointType("Gopher"),
	Things:  EndPointType("Things"),
	Thing:   EndPointType("Thing"),
}

type endPointTypes struct {
	Gophers EndPointType
	Gopher  EndPointType
	Things  EndPointType
	Thing   EndPointType
}

func (c EndPointType) String() string {
	return "pkg.metrics.endpoints." + string(c)
}

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
	command *prometheus.CounterVec
}

type serviceMetric struct {
	transport *prometheus.CounterVec
}

func NewMetrics() (m *Metrics) {
	m = new(Metrics)

	m.serverInit()
	m.serviceInit()
	m.clientInit()

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
	m.service.transport = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: `gophercli_service_transport_total`,
			Help: "How many services calls",
		}, []string{"method", "endpoint", "status"},
	)
}
func (m *Metrics) clientInit() {
	m.client.command = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: `gophercli_client_total`,
			Help: "How many CLI client / command calls were made",
		}, []string{"method", "endpoint"},
	)
}

func (m *Metrics) Marshal(filename string) {
	prometheus.WriteToTextfile(filename, prometheus.DefaultGatherer)
}

func (m *Metrics) ServerInc(endPoint EndPointType, method serviceMethodType) {
	if m.server.endPoint == nil {
		return
	}
	labels := prometheus.Labels{"endpoint": endPoint.String(), "method": method.String()}
	m.server.endPoint.With(labels).Inc()
}

func (m *Metrics) DBInc(endPoint EndPointType, method dbMethodType) {
	if m.server.db == nil {
		return
	}
	labels := prometheus.Labels{"endpoint": endPoint.String(), "method": method.String()}
	m.server.db.With(labels).Inc()
}

func (m *Metrics) CacheInc(endPoint EndPointType, method cacheMethodType) {
	if m.server.cache == nil {
		return
	}
	labels := prometheus.Labels{"endpoint": endPoint.String(), "method": method.String()}
	m.server.cache.With(labels).Inc()
}

func (m *Metrics) TransportInc(endPoint EndPointType, method transportMethodType, status int) {
	if m.service.transport == nil {
		return
	}
	labels := prometheus.Labels{"endpoint": endPoint.String(), "method": method.String(), "status": fmt.Sprintf("%d", status)}
	m.service.transport.With(labels).Inc()
}

func (m *Metrics) ClientInc(endPoint string, method serviceMethodType) {
	if m.client.command == nil {
		return
	}
	labels := prometheus.Labels{"endpoint": endPoint, "method": method.String()}
	m.client.command.With(labels).Inc()
}

func DumpMetrics(file string) {
	prometheus.WriteToTextfile(file, prometheus.DefaultGatherer)
}
