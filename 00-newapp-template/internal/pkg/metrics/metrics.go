package metrics

import (
	"00-newapp-template/pkg/acme"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type Metrics struct {
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

func NewMetrics() (m *Metrics) {
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

var Methods = methodTypes{
	Service: serviceTypes{
		Get:    serviceMethodType("Get"),
		Update: serviceMethodType("Update"),
		Add:    serviceMethodType("Add"),
		Delete: serviceMethodType("Delete"),
	},
	DB: dbTypes{
		Delete: dbMethodType("Delete"),
		Update: dbMethodType("Update"),
		Read:   dbMethodType("Read"),
		Insert: dbMethodType("Insert"),
	},
	Cache: cacheTypes{
		Hit:        cacheMethodType("Hit"),
		Miss:       cacheMethodType("Miss"),
		Invalidate: cacheMethodType("Invalidate"),
		Store:      cacheMethodType("Store"),
	},
}

type methodTypes struct {
	Service   serviceTypes
	DB        dbTypes
	Cache     cacheTypes
	Transport transportTypes
}

type serviceMethodType string
type serviceTypes struct {
	Get    serviceMethodType
	Update serviceMethodType
	Add    serviceMethodType
	Delete serviceMethodType
}

func (c serviceMethodType) String() string {
	return "pkg.metric.service." + string(c)
}

type dbMethodType string
type dbTypes struct {
	Read   dbMethodType
	Update dbMethodType
	Insert dbMethodType
	Delete dbMethodType
}

func (c dbMethodType) String() string {
	return "pkg.metric.db." + string(c)
}

type cacheMethodType string
type cacheTypes struct {
	Hit        cacheMethodType
	Miss       cacheMethodType
	Store      cacheMethodType
	Invalidate cacheMethodType
}

func (c cacheMethodType) String() string {
	return "pkg.metric.cache." + string(c)
}

type transportMethodType string
type transportTypes struct {
	Get    transportMethodType
	Put    transportMethodType
	Post   transportMethodType
	Delete transportMethodType
	Head   transportMethodType
}

func (c transportMethodType) String() string {
	return "pkg.metric.transport." + string(c)
}

func (m *Metrics) ServiceInc(endPoint acme.ServiceEndPoint, method serviceMethodType) {
	labels := prometheus.Labels{"endpoint": endPoint.String(), "method": method.String()}
	m.server.EndPoint.With(labels).Inc()
}

func (m *Metrics) DBInc(endPoint acme.ServiceEndPoint, method dbMethodType) {
	labels := prometheus.Labels{"endpoint": endPoint.String(), "method": method.String()}
	m.server.DB.With(labels).Inc()
}

func (m *Metrics) CacheInc(endPoint acme.ServiceEndPoint, method cacheMethodType) {
	labels := prometheus.Labels{"endpoint": endPoint.String(), "method": method.String()}
	m.server.Cache.With(labels).Inc()
}

// GET/POST/PUT/DELETE
func (m *Metrics) TransportInc(endPoint acme.ServiceEndPoint, method transportMethodType) {
	labels := prometheus.Labels{"endpoint": endPoint.String(), "method": method.String()}
	m.service.Transport.With(labels).Inc()
}
