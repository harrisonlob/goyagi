package metrics

import (
	"fmt"
	"time"

	"github.com/DataDog/datadog-go/statsd"
	"github.com/harrisonlob/goyagi/pkg/config"
	"github.com/labstack/echo"
)

type statsdClient interface {
	Histogram(name string, value float64, tags []string, rate float64) error
	Count(name string, value int64, tags []string, rate float64) error
}

type Metrics struct {
	client statsdClient
}

type Timer struct {
	name    string
	metrics *Metrics
	begin   time.Time
	tags    []string
}

const namespace = "goyagi."

func New(cfg config.Config) (Metrics, error) {
	address := fmt.Sprintf("%s:%d", cfg.StatsdHost, cfg.StatsdPort)

	client, err := statsd.New(address)
	if err != nil {
		return Metrics{}, err
	}

	client.Namespace = namespace
	client.Tags = []string{
		fmt.Sprintf("environment:%s", cfg.Environment),
	}

	return Metrics{client}, nil
}

func (m *Metrics) Count(name string, count int64, tags ...string) {
	m.client.Count(name, count, tags, 1) // nolint:gosec
}

func (m *Metrics) Histogram(name string, value float64, tags ...string) {
	m.client.Histogram(name, value, tags, 1) // nolint:gosec
}

func (m *Metrics) NewTimer(name string, tags ...string) Timer {
	return Timer{
		begin:   time.Now(),
		metrics: m,
		name:    name,
		tags:    tags,
	}
}

func (t *Timer) End(additionalTags ...string) float64 {
	duration := time.Since(t.begin)
	durationInMS := float64(duration / time.Millisecond)

	t.tags = append(t.tags, additionalTags...)

	t.metrics.Histogram(t.name, durationInMS, t.tags...)

	return durationInMS
}

// Middleware returns an Echo middleware function that begins a timer before a
// request is handled and ends afterwards.
func Middleware(m Metrics) func(next echo.HandlerFunc) echo.HandlerFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			methodTag := fmt.Sprintf("method:%s", c.Request().Method)

			t := m.NewTimer("http.request", methodTag)

			if err := next(c); err != nil {
				c.Error(err)
			}

			statusCodeTag := fmt.Sprintf("status_code:%d", c.Response().Status)
			pathTag := fmt.Sprintf("path:%s", c.Path())

			t.End(statusCodeTag, pathTag)

			return nil
		}
	}
}
