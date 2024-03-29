package telemetry

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	api "go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	reportInterval = time.Second * 15
	reportTimeout  = time.Second * 3
)

type Option func(*telemetryAgent)

func ReportInterval(v time.Duration) Option {
	return func(t *telemetryAgent) {
		t.reportInterval = v
	}
}

func ReportTimeout(v time.Duration) Option {
	return func(t *telemetryAgent) {
		t.reportTimeout = v
	}
}

func GRPCConn(v *grpc.ClientConn) Option {
	return func(t *telemetryAgent) {
		t.conn = v
	}
}

type telemetryAgent struct {
	meter          api.Meter
	conn           *grpc.ClientConn
	reportInterval time.Duration
	reportTimeout  time.Duration
}

func NewTelemetryAgent(addr string, opts ...Option) (*telemetryAgent, error) {
	if addr == "" {
		return nil, errors.Errorf("missing metrics rpc agent address")
	}

	m := &telemetryAgent{
		reportInterval: reportInterval,
		reportTimeout:  reportTimeout,
	}
	for _, f := range opts {
		f(m)
	}

	if m.conn == nil {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()
		conn, err := grpc.DialContext(ctx, addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			return nil, err
		}
		m.conn = conn
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	exporter, err := otlpmetricgrpc.New(ctx, otlpmetricgrpc.WithGRPCConn(m.conn))
	if err != nil {
		return nil, err
	}

	res, err := resource.New(context.Background(), resource.WithAttributes())
	if err != nil {
		return nil, err
	}

	provider := metric.NewMeterProvider(
		metric.WithReader(metric.NewPeriodicReader(exporter,
			metric.WithTimeout(m.reportTimeout),
			metric.WithInterval(m.reportInterval),
		)), metric.WithResource(res))
	otel.SetMeterProvider(provider)

	m.meter = provider.Meter("_exporter_")
	return m, nil
}

var (
	ErrInvalidMetric = errors.New("metrics meter is invalid")
)

// Int64Counter returns a new Int64Counter instrument identified by name
// and configured with options. The instrument is used to synchronously
// record increasing int64 measurements during a computational operation.
func (m *telemetryAgent) Int64Counter(name string, options ...api.Int64CounterOption) (api.Int64Counter, error) {
	if m == nil || m.meter == nil {
		return nil, ErrInvalidMetric
	}
	return m.meter.Int64Counter(name, options...)
}

// Int64UpDownCounter returns a new Int64UpDownCounter instrument
// identified by name and configured with options. The instrument is used
// to synchronously record int64 measurements during a computational
// operation.
func (m *telemetryAgent) Int64UpDownCounter(name string, options ...api.Int64UpDownCounterOption) (api.Int64UpDownCounter, error) {
	if m == nil || m.meter == nil {
		return nil, ErrInvalidMetric
	}
	return m.meter.Int64UpDownCounter(name, options...)
}

// Int64Histogram returns a new Int64Histogram instrument identified by
// name and configured with options. The instrument is used to
// synchronously record the distribution of int64 measurements during a
// computational operation.
// metric.WithExplicitBucketBoundaries to exclipt set buckets
func (m *telemetryAgent) Int64Histogram(name string, options ...api.Int64HistogramOption) (api.Int64Histogram, error) {
	if m == nil || m.meter == nil {
		return nil, ErrInvalidMetric
	}
	return m.meter.Int64Histogram(name, options...)
}

// Float64Counter returns a new Float64Counter instrument identified by
// name and configured with options. The instrument is used to
// synchronously record increasing float64 measurements during a
// computational operation.
func (m *telemetryAgent) Float64Counter(name string, options ...api.Float64CounterOption) (api.Float64Counter, error) {
	if m == nil || m.meter == nil {
		return nil, ErrInvalidMetric
	}
	return m.meter.Float64Counter(name, options...)
}

// Float64UpDownCounter returns a new Float64UpDownCounter instrument
// identified by name and configured with options. The instrument is used
// to synchronously record float64 measurements during a computational
// operation.
func (m *telemetryAgent) Float64UpDownCounter(name string, options ...api.Float64UpDownCounterOption) (api.Float64UpDownCounter, error) {
	if m == nil || m.meter == nil {
		return nil, ErrInvalidMetric
	}
	return m.meter.Float64UpDownCounter(name, options...)
}

// Float64Histogram returns a new Float64Histogram instrument identified by
// name and configured with options. The instrument is used to
// synchronously record the distribution of float64 measurements during a
// computational operation.
func (m *telemetryAgent) Float64Histogram(name string, options ...api.Float64HistogramOption) (api.Float64Histogram, error) {
	if m == nil || m.meter == nil {
		return nil, ErrInvalidMetric
	}
	return m.meter.Float64Histogram(name, options...)
}

// RegisterCallback registers f to be called during the collection of a
// measurement cycle.
//
// If Unregister of the returned Registration is called, f needs to be
// unregistered and not called during collection.
//
// The instruments f is registered with are the only instruments that f may
// observe values for.
//
// If no instruments are passed, f should not be registered nor called
// during collection.
//
// The function f needs to be concurrent safe.
func (m *telemetryAgent) RegisterInt64CounterCallback(f Int64Callback, name string, options ...api.Int64ObservableCounterOption) (api.Registration, error) {
	if m == nil || m.meter == nil {
		return nil, ErrInvalidMetric
	}

	c, err := m.meter.Int64ObservableCounter(name, options...)
	if err != nil {
		return nil, err
	}
	return m.meter.RegisterCallback(func(ctx context.Context, o api.Observer) error {
		v, err := f()
		o.ObserveInt64(c, v.Value, v.Options...)
		return err
	}, c)
}

func (m *telemetryAgent) RegisterInt64UpDownCounterCallback(f Int64Callback, name string, options ...api.Int64ObservableUpDownCounterOption) (api.Registration, error) {
	if m == nil || m.meter == nil {
		return nil, ErrInvalidMetric
	}

	c, err := m.meter.Int64ObservableUpDownCounter(name, options...)
	if err != nil {
		return nil, err
	}
	return m.meter.RegisterCallback(func(ctx context.Context, o api.Observer) error {
		v, err := f()
		o.ObserveInt64(c, v.Value, v.Options...)
		return err
	}, c)
}

func (m *telemetryAgent) RegisterInt64GaugeCallback(f Int64Callback, name string, options ...api.Int64ObservableGaugeOption) (api.Registration, error) {
	if m == nil || m.meter == nil {
		return nil, ErrInvalidMetric
	}

	g, err := m.meter.Int64ObservableGauge(name, options...)
	if err != nil {
		return nil, err
	}
	return m.meter.RegisterCallback(func(ctx context.Context, o api.Observer) error {
		v, err := f()
		o.ObserveInt64(g, v.Value, v.Options...)
		return err
	}, g)
}

func (m *telemetryAgent) RegisterFloat64CounterCallback(f Float64Callback, name string, options ...api.Float64ObservableCounterOption) (api.Registration, error) {
	if m == nil || m.meter == nil {
		return nil, ErrInvalidMetric
	}

	c, err := m.meter.Float64ObservableCounter(name, options...)
	if err != nil {
		return nil, err
	}
	return m.meter.RegisterCallback(func(ctx context.Context, o api.Observer) error {
		v, err := f()
		o.ObserveFloat64(c, v.Value, v.Options...)
		return err
	}, c)
}

func (m *telemetryAgent) RegisterFloat64UpDownCounterCallback(f Float64Callback, name string, options ...api.Float64ObservableUpDownCounterOption) (api.Registration, error) {
	if m == nil || m.meter == nil {
		return nil, ErrInvalidMetric
	}

	c, err := m.meter.Float64ObservableUpDownCounter(name, options...)
	if err != nil {
		return nil, err
	}
	return m.meter.RegisterCallback(func(ctx context.Context, o api.Observer) error {
		v, err := f()
		o.ObserveFloat64(c, v.Value, v.Options...)
		return err
	}, c)
}

func (m *telemetryAgent) RegisterFloat64GaugeCallback(f Float64Callback, name string, options ...api.Float64ObservableGaugeOption) (api.Registration, error) {
	if m == nil || m.meter == nil {
		return nil, ErrInvalidMetric
	}

	g, err := m.meter.Float64ObservableGauge(name, options...)
	if err != nil {
		return nil, err
	}
	return m.meter.RegisterCallback(func(ctx context.Context, o api.Observer) error {
		v, err := f()
		o.ObserveFloat64(g, v.Value, v.Options...)
		return err
	}, g)
}

type Int64Callback func() (ObserverInt64, error)
type ObserverInt64 struct {
	Value   int64
	Options []api.ObserveOption
}

type Float64Callback func() (ObserverFloat64, error)
type ObserverFloat64 struct {
	Value   float64
	Options []api.ObserveOption
}
