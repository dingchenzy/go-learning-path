package controllers

import (
	"database/sql"

	"github.com/prometheus/client_golang/prometheus"
)

type ThreadsCollector struct {
	*mysqlSelcetor
	desc *prometheus.Desc
}

func NewThreadCoThreadsCollector(db *sql.DB) *ThreadsCollector {
	return &ThreadsCollector{
		mysqlSelcetor: newMysqlSelcector(db),
		desc:          prometheus.NewDesc("mysql_thread_count", "current thread count", nil, nil),
	}
}

func (t *ThreadsCollector) Describe(descs chan<- *prometheus.Desc) {
	descs <- t.desc
}

func (t *ThreadsCollector) Collect(metrics chan<- prometheus.Metric) {
	value := t.mysqlSelcetor.showStatus("Threads_connected")
	metrics <- prometheus.MustNewConstMetric(t.desc, prometheus.GaugeValue, value)
}
