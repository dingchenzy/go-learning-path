# 开发redis_exporter

## 监控指标

监控 info 内指标，案例中监控的内容为。

```sh
connected_clients		#客户端连接数
使用 Ping().Result()	#方法返回的关键字判断 redis 是否存活
```

## 项目目录结构

```bash
.
├── controllera			# 存放处理器
│   ├── base.go		# 查询功能
│   ├── clientconn.go	# 客户端连接数，形成 collector 结构体
│   └── up.go	# 存活检测
├── go.mod
├── go.sum
├── main.go		# 主文件，提供注册功能
├── redis_exporter
└── vendor
    ├── github.com
    ├── golang.org
    ├── google.golang.org
    └── modules.txt
```

## 代码主体

### main.go

```go
package main

import (
	"log"
	"net/http"

	"github.com/exporter/redis_exporter/controllera"
	"github.com/go-redis/redis"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	addr     = "10.0.0.98:6379"
	httpaddr = "0.0.0.0:9998"
)

func NewClient() *redis.Client {
	client := redis.NewClient(
		&redis.Options{
			Network: "tcp",
			Addr:    addr,
		},
	)
	pong, err := client.Ping().Result()
	if err != nil || pong != "PONG" {
		log.Fatal("[error]connection redis is failed：", err)
	}
	return client
}

func main() {
	client := NewClient()

	prometheus.Register(controllera.NewClientConnCollector(client))
	prometheus.Register(controllera.NewUpRedisCollector(client))

	http.Handle("/metrics/", promhttp.Handler())

	http.ListenAndServe(httpaddr, nil)
}
```

### base.go

```go
package controllera

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/go-redis/redis"
)

type baseSelect struct {
	redisconn *redis.Client
}

func newbaseSelect(conredis *redis.Client) *baseSelect {
	return &baseSelect{
		redisconn: conredis,
	}
}

func (b *baseSelect) redisele(grepvalue string) (float64, error) {
	var str string

	name := b.redisconn.Info()

	err := name.Scan(&str)
	if err != nil {
		log.Print("[error]read redis info is error :", err)
	}

	regex, err := regexp.Compile("\r\n")
	if err != nil {
		log.Print("[error]read redis info regexp the is :", err)
	}

	str1 := regex.Split(str, -1)

	for _, v := range str1 {
		if strings.Contains(v, grepvalue) {
			value := strings.Split(v, ":")
			if num, err := strconv.Atoi(value[1]); err == nil {
				return float64(num), nil
			} else {
				log.Print("[error]string strconv is error : ", err)
			}
		}
	}
	return 0, fmt.Errorf("grepvalue is not found")
}
```

### clientconn.go

```go
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
```

### up.go

```go
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
```

