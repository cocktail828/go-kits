package prome

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

type promeServer struct {
	lock         sync.RWMutex
	registry     *prometheus.Registry
	collectorMap map[string]prometheus.Collector
	constLables  map[string]string
}

type Option func(*promeServer)

func WithProcessCollector() Option {
	return func(ms *promeServer) {
		ms.registry.MustRegister(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))
	}
}

func WithGoCollector() Option {
	return func(ms *promeServer) {
		ms.registry.MustRegister(collectors.NewGoCollector())
	}
}

func WithConstLables(lables map[string]string) Option {
	return func(ms *promeServer) {
		for k, v := range lables {
			ms.constLables[k] = v
		}
	}
}

func NewMetricsServer(r *prometheus.Registry, opts ...Option) *promeServer {
	if r == nil {
		r = prometheus.NewRegistry()
	}

	srv := &promeServer{
		registry:     r,
		collectorMap: make(map[string]prometheus.Collector),
		constLables:  make(map[string]string),
	}

	for _, f := range opts {
		f(srv)
	}

	return srv
}

func (ms *promeServer) setCollector(opt CollectorOpt, collector prometheus.Collector) {
	ms.lock.Lock()
	defer ms.lock.Unlock()

	ms.registry.MustRegister(collector)
	ms.collectorMap[opt.String()] = collector
}

func (ms *promeServer) UnregisterByOpts(opt CollectorOpt) {
	ms.lock.RLock()
	defer ms.lock.RUnlock()
	collector, ok := ms.collectorMap[opt.String()]
	if ok {
		ms.UnregisterByCollector(collector)
	}
}

func (ms *promeServer) UnregisterByCollector(collector prometheus.Collector) {
	ms.lock.Lock()
	defer ms.lock.Unlock()
	ms.registry.Unregister(collector)
}

func (ms *promeServer) RegisterGauge(opt CollectorOpt) prometheus.Gauge {
	collector := prometheus.NewGauge(prometheus.GaugeOpts(opt))
	ms.setCollector(opt, collector)
	return collector
}

func (ms *promeServer) RegisterCounter(opt CollectorOpt) prometheus.Counter {
	collector := prometheus.NewCounter(prometheus.CounterOpts(opt))
	ms.setCollector(opt, collector)
	return collector
}

func (ms *promeServer) RegisterHistogram(opt CollectorOpt, buckets []float64) prometheus.Histogram {
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

func (ms *promeServer) RegisterSummary(opt CollectorOpt, objectives map[float64]float64) prometheus.Summary {
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

func (ms *promeServer) RegisterGaugeVec(opt CollectorOpt, labels []string) *prometheus.GaugeVec {
	collector := prometheus.NewGaugeVec(prometheus.GaugeOpts(opt), labels)
	ms.setCollector(opt, collector)
	return collector
}

func (ms *promeServer) RegisterCounterVec(opt CollectorOpt, labels []string) *prometheus.CounterVec {
	collector := prometheus.NewCounterVec(prometheus.CounterOpts(opt), labels)
	ms.setCollector(opt, collector)
	return collector
}

func (ms *promeServer) RegisterHistogramVec(opt CollectorOpt, buckets []float64, labels []string) *prometheus.HistogramVec {
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

func (ms *promeServer) RegisterSummaryVec(opt CollectorOpt, objectives map[float64]float64, labels []string) *prometheus.SummaryVec {
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

func (ms *promeServer) Run(addr string) {
	// Serve the default Prometheus promeServer registry over HTTP on /promeServer.
	http.Handle("/metrics", promhttp.HandlerFor(ms.registry, promhttp.HandlerOpts{Registry: ms.registry}))
	if err := http.ListenAndServe(addr, nil); err != nil {
		fmt.Println(err)
	}
}
