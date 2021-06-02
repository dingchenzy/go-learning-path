package controllera

import (
	"log"

	"github.com/go-redis/redis"
	"github.com/prometheus/client_golang/prometheus"
)

type UpRedisCollector struct {
	*baseSelect
	desc *prometheus.Desc
}

func NewUpRedisCollector(rediscli *redis.Client) *UpRedisCollector {
	return &UpRedisCollector{
		baseSelect: &baseSelect{
			redisconn: rediscli,
		},
		desc: prometheus.NewDesc("redis_life", "this is redis life", nil, nil),
	}
}

func (u *UpRedisCollector) Describe(descs chan<- *prometheus.Desc) {
	descs <- u.desc
}

func (u *UpRedisCollector) Collect(metrics chan<- prometheus.Metric) {
	var value int = 1
	if v, err := u.baseSelect.redisconn.Ping().Result(); v != "PONG" || err != nil {
		log.Print(err)
		value = 0
	}
	metrics <- prometheus.MustNewConstMetric(u.desc, prometheus.GaugeValue, float64(value))
}
