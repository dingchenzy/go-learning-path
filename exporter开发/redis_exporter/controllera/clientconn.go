package controllera

import (
	"log"

	"github.com/go-redis/redis"
	"github.com/prometheus/client_golang/prometheus"
)

type ClientConnCollector struct {
	*baseSelect
	desc *prometheus.Desc
}

func NewClientConnCollector(rediscli *redis.Client) *ClientConnCollector {
	return &ClientConnCollector{
		baseSelect: &baseSelect{
			redisconn: rediscli,
		},
		desc: prometheus.NewDesc("client_conn_count", "this is connection count", nil, nil),
	}
}

func (c *ClientConnCollector) Describe(descs chan<- *prometheus.Desc) {
	descs <- c.desc
}

func (c *ClientConnCollector) Collect(metrics chan<- prometheus.Metric) {
	value, err := c.baseSelect.redisele("connected_clients")
	if err != nil {
		log.Print("[error]get baseelect readisele value is error : ", err)
	}
	metrics <- prometheus.MustNewConstMetric(c.desc, prometheus.GaugeValue, value)
}
