package controllers

import (
	"database/sql"

	"github.com/prometheus/client_golang/prometheus"
)

type LongtimeCollector struct {
	*mysqlSelcetor
	desc *prometheus.Desc
}

func NewLongtimeLongtimeCollector(sqldb *sql.DB) *LongtimeCollector {
	return &LongtimeCollector{
		mysqlSelcetor: &mysqlSelcetor{sqldb},
		desc:          prometheus.NewDesc("mysql_select_long_time", "this is select timeout", nil, nil),
	}
}

func (l *LongtimeCollector) Describe(desc chan<- *prometheus.Desc) {
	desc <- l.desc
}

func (l *LongtimeCollector) Collect(metric chan<- prometheus.Metric) {
	metric <- prometheus.MustNewConstMetric(l.desc, prometheus.CounterValue, l.mysqlSelcetor.showStatus("Slow_queries"))
}
