package metrics

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// total number of tasks created
var TasksCreatedTotal = prometheus.NewCounter(prometheus.CounterOpts{
	Name: "tasks_created_total",
	Help: "Total number of tasks created",
})

// creating task
var TaskCreateDuration = prometheus.NewHistogram(prometheus.HistogramOpts{
	Name:    "task_create_duration_seconds",
	Help:    "Time spent creating a task",
	Buckets: []float64{0.1, 0.5, 1, 5, 15},
})

// executing task
var TaskExecutionDuration = prometheus.NewHistogram(prometheus.HistogramOpts{
	Name:    "task_execution_duration_seconds",
	Help:    "Time spent executing a task",
	Buckets: []float64{0.1, 0.5, 1, 5, 15},
})

// errors during task creation
var TaskCreateErrorsTotal = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "task_create_errors_total",
		Help: "Total number of errors during task creation",
	},
	[]string{"error"},
)

// errors during task execution
var TaskExecutionErrorsTotal = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "task_execution_errors_total",
		Help: "Total number of errors during task execution",
	},
	[]string{"error"},
)

func init() {
	prometheus.MustRegister(
		TasksCreatedTotal,
		TaskCreateDuration,
		TaskExecutionDuration,
		TaskCreateErrorsTotal,
		TaskExecutionErrorsTotal,
	)
}

func StartMetricsServer(port string) error {
	http.Handle("/metrics", promhttp.Handler())
	return http.ListenAndServe(port, nil)
}
