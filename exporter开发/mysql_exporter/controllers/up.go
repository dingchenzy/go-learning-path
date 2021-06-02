package controllers

import (
	"database/sql"

	"github.com/prometheus/client_golang/prometheus"
)

type UpController struct {
	*mysqlSelcetor
	desc *prometheus.Desc
}

// 构建函数，并封装 newMysqlSelector 中的查询方法
func NewUpCollertor(sqldb *sql.DB) *UpController {
	return &UpController{
		newMysqlSelcector(sqldb),
		prometheus.NewDesc("mysql_up", "mysql up or down", nil, nil),
	}
}

// desc 中集成着 label 的名称，指标的名称，help 的内容等
func (Up *UpController) Describe(descs chan<- *prometheus.Desc) {
	descs <- Up.desc
}

func (Up *UpController) Collect(metrics chan<- prometheus.Metric) {
	var up float64 = 1
	if err := Up.mysqlSelcetor.mysqlconnection.Ping(); err != nil {
		up = 0
	}
	// 生成 metrics
	// metrics 内部集成着 desc 以及指标的类型，label 的 value 等内容
	metrics <- prometheus.MustNewConstMetric(Up.desc, prometheus.GaugeValue, up)
}
