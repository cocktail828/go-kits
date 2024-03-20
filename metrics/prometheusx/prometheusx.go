package prometheusx

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type CollectorOpt prometheus.Opts

func (c CollectorOpt) String() string {
	return fmt.Sprintf("%s_%s_%s", c.Namespace, c.Subsystem, c.Name)
}

type prometheusServer struct {
	lock         sync.RWMutex
	registry     *prometheus.Registry
	collectorMap map[string]prometheus.Collector
	constLables  map[string]string
}

type Option func(*prometheusServer)

func WithProcessCollector() Option {
	return func(ms *prometheusServer) {
		ms.registry.MustRegister(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))
	}
}

func WithGoCollector() Option {
	return func(ms *prometheusServer) {
		ms.registry.MustRegister(collectors.NewGoCollector())
	}
}

func WithConstLables(lables map[string]string) Option {
	return func(ms *prometheusServer) {
		for k, v := range lables {
			ms.constLables[k] = v
		}
	}
}

func NewMetricsServer(r *prometheus.Registry, opts ...Option) *prometheusServer {
	if r == nil {
		r = prometheus.NewRegistry()
	}

	srv := &prometheusServer{
		registry:     r,
		collectorMap: make(map[string]prometheus.Collector),
		constLables:  make(map[string]string),
	}

	for _, f := range opts {
		f(srv)
	}

	return srv
}

func (ms *prometheusServer) setCollector(opt CollectorOpt, collector prometheus.Collector) {
	ms.lock.Lock()
	defer ms.lock.Unlock()

	ms.registry.MustRegister(collector)
	ms.collectorMap[opt.String()] = collector
}

func (ms *prometheusServer) UnregisterByOpts(opt CollectorOpt) {
	ms.lock.RLock()
	defer ms.lock.RUnlock()
	collector, ok := ms.collectorMap[opt.String()]
	if ok {
		ms.UnregisterByCollector(collector)
	}
}

func (ms *prometheusServer) UnregisterByCollector(collector prometheus.Collector) {
	ms.lock.Lock()
	defer ms.lock.Unlock()
	ms.registry.Unregister(collector)
}

func (ms *prometheusServer) RegisterGauge(opt CollectorOpt) prometheus.Gauge {
	collector := prometheus.NewGauge(prometheus.GaugeOpts(opt))
	ms.setCollector(opt, collector)
	return collector
}

func (ms *prometheusServer) RegisterCounter(opt CollectorOpt) prometheus.Counter {
	collector := prometheus.NewCounter(prometheus.CounterOpts(opt))
	ms.setCollector(opt, collector)
	return collector
}

func (ms *prometheusServer) RegisterHistogram(opt CollectorOpt, buckets []float64) prometheus.Histogram {
	collector := prometheus.NewHistogram(prometheus.HistogramOpts{
		Namespace:   opt.Namespace,
		Subsystem:   opt.Subsystem,
		Name:        opt.Name,
		Help:        opt.Help,
		ConstLabels: opt.ConstLabels,
		Buckets:     buckets,
	})
	ms.setCollector(opt, collector)
	return collector
}

func (ms *prometheusServer) RegisterSummary(opt CollectorOpt, objectives map[float64]float64) prometheus.Summary {
	collector := prometheus.NewSummary(prometheus.SummaryOpts{
		Namespace:   opt.Namespace,
		Subsystem:   opt.Subsystem,
		Name:        opt.Name,
		Help:        opt.Help,
		ConstLabels: opt.ConstLabels,
		Objectives:  objectives,
	})
	ms.setCollector(opt, collector)
	return collector
}

func (ms *prometheusServer) RegisterGaugeVec(opt CollectorOpt, labels []string) *prometheus.GaugeVec {
	collector := prometheus.NewGaugeVec(prometheus.GaugeOpts(opt), labels)
	ms.setCollector(opt, collector)
	return collector
}

func (ms *prometheusServer) RegisterCounterVec(opt CollectorOpt, labels []string) *prometheus.CounterVec {
	collector := prometheus.NewCounterVec(prometheus.CounterOpts(opt), labels)
	ms.setCollector(opt, collector)
	return collector
}

func (ms *prometheusServer) RegisterHistogramVec(opt CollectorOpt, buckets []float64, labels []string) *prometheus.HistogramVec {
	collector := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace:   opt.Namespace,
		Subsystem:   opt.Subsystem,
		Name:        opt.Name,
		Help:        opt.Help,
		ConstLabels: opt.ConstLabels,
		Buckets:     buckets,
	}, labels)
	ms.setCollector(opt, collector)
	return collector
}

func (ms *prometheusServer) RegisterSummaryVec(opt CollectorOpt, objectives map[float64]float64, labels []string) *prometheus.SummaryVec {
	collector := prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Namespace:   opt.Namespace,
		Subsystem:   opt.Subsystem,
		Name:        opt.Name,
		Help:        opt.Help,
		ConstLabels: opt.ConstLabels,
		Objectives:  objectives,
	}, labels)
	ms.setCollector(opt, collector)
	return collector
}

func (ms *prometheusServer) Run(addr string) {
	// Serve the default Prometheus prometheusServer registry over HTTP on /prometheusServer.
	http.Handle("/metrics", promhttp.HandlerFor(ms.registry, promhttp.HandlerOpts{Registry: ms.registry}))
	if err := http.ListenAndServe(addr, nil); err != nil {
		fmt.Println(err)
	}
}
