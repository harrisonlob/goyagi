package metrics

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/harrisonlob/goyagi/pkg/config"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testCount = int64(1)
const testDuration = float64(50)
const testRate = float64(1)
const testMetric = "test_metric"
const testTag = "foo:bar"

var testTags = []string{testTag}

type mockClient struct {
	t        *testing.T
	name     string
	count    int64
	duration float64
	tags     []string
	rate     float64
}

func (m *mockClient) Count(name string, count int64, tags []string, rate float64) error {
	m.name = testMetric
	m.count = testCount
	m.tags = testTags
	m.rate = testRate

	return errors.New("test error")
}

func (m *mockClient) Histogram(name string, duration float64, tags []string, rate float64) error {
	m.name = testMetric
	m.duration = testDuration
	m.tags = testTags
	m.rate = testRate

	return errors.New("test error")
}

func newMockedClient(t *testing.T, cfg config.Config) Metrics {
	metrics, err := New(cfg)
	require.NoError(t, err)
	require.NotNil(t, metrics)

	metrics.client = &mockClient{t, "", 0, 0, []string{}, 0}

	return metrics
}

func TestCount(t *testing.T) {
	cfg := config.New()

	t.Run("calls Datadog Count function and ignores error", func(tt *testing.T) {
		metrics := newMockedClient(t, cfg)

		metrics.Count(testMetric, testCount, testTags...)

		mc, ok := metrics.client.(*mockClient)
		require.True(t, ok, "unexpected error during type assertion")

		assert.Equal(t, testMetric, mc.name, "inconsistent metric name")
		assert.Equal(t, testCount, mc.count, "inconsistent metric count")
		assert.Equal(t, testTags, mc.tags, "inconsistent tags")
		assert.Equal(t, testRate, mc.rate, "inconsistent rate")
	})
}

func TestHistogram(t *testing.T) {
	cfg := config.New()

	t.Run("calls Datadog Histogram function and ignores error", func(tt *testing.T) {
		metrics := newMockedClient(t, cfg)

		metrics.Histogram(testMetric, testDuration, testTags...)

		mc, ok := metrics.client.(*mockClient)
		require.True(t, ok, "unexpected error during type assertion")

		assert.Equal(t, testMetric, mc.name, "inconsistent metric name")
		assert.Equal(t, testDuration, mc.duration, "inconsistent duration")
		assert.Equal(t, testTags, mc.tags, "inconsistent tags")
		assert.Equal(t, testRate, mc.rate, "inconsistent rate")
	})
}

func TestTimer(t *testing.T) {
	cfg := config.New()

	t.Run("calls histogram with correct duration", func(tt *testing.T) {
		metrics := newMockedClient(t, cfg)

		timer := metrics.NewTimer(testMetric, testTag)
		require.NotNil(t, timer)

		time.Sleep(time.Duration(testDuration) * time.Millisecond)
		timer.End()

		mc, ok := metrics.client.(*mockClient)
		require.True(t, ok, "unexpected error during type assertion")

		assert.Equal(t, testMetric, mc.name, "inconsistent metric name")
		assert.True(t, testDuration <= mc.duration, "incorrect duration")
		assert.Equal(t, testTags, mc.tags, "inconsistent tags")
		assert.Equal(t, testRate, mc.rate, "inconsistent rate")
	})
}

func TestMiddleware(t *testing.T) {
	cfg := config.New()

	t.Run("sends request duration through Datadog client", func(tt *testing.T) {
		metrics := newMockedClient(t, cfg)

		e := echo.New()
		e.Use(Middleware(metrics))

		e.GET("/", func(c echo.Context) error {
			time.Sleep(time.Duration(testDuration) * time.Millisecond)
			return nil
		})

		req, err := http.NewRequest("GET", "/", nil)
		require.NoError(t, err)

		rr := httptest.NewRecorder()
		e.ServeHTTP(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)

		mc, ok := metrics.client.(*mockClient)
		require.True(t, ok, "unexpected error during type assertion")

		assert.True(t, testDuration <= mc.duration, "incorrect duration")
	})
}
