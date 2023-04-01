package main

import (
	"log"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	counter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "myapp_counter",
		Help: "myapp_counter_help",
	})

	counterVec = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "myapp_counter_vec",
		Help: "myapp_counter_vec_help",
	}, []string{"action"})

	histogramVec = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "myapp_histogram_vec",
		Help:    "myapp_histogram_vec_help",
		Buckets: []float64{10, 20, 50, 100},
	}, []string{"action"})
)

func DoMetricHttpRequest(startTime time.Time) {
	d := time.Now().Sub(startTime)
	ms := d.Milliseconds()
	lables := map[string]string{"action": "http_request"} // 可以传入具体的action
	histogramVec.With(lables).Observe(float64(ms))
	counterVec.With(lables).Inc()
}

func DoMetricProbeCount(probe string) {
	lables := map[string]string{"action": probe}
	counterVec.With(lables).Inc()
}

var (
	promHandler = promhttp.Handler()
)

type MetricHandler struct {
}

func (p *MetricHandler) ServeHTTP(writer http.ResponseWriter, req *http.Request) {
	log.Printf("metric-ServeHTTP")
	promHandler.ServeHTTP(writer, req)
}

func InitMetric() {
	go func() {
		if err := http.ListenAndServe(":17181", &MetricHandler{}); err != nil {
			log.Printf("metric ListenAndServe fail, err=%v", err)
		}
	}()

	log.Printf("InitMetric ok")
}

func DeInitMetric() {
	log.Printf("DeInitMetric ok")
}
