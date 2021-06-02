package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/exporter/mysql_exporter/controllers"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	_ "github.com/go-sql-driver/mysql"
)

const (
	dsn    = "localuser:123.com@tcp(localhost:3306)/mysql?charset=utf8mb4&loc=PRC&parseTime=true"
	driver = "mysql"
	addr   = ":9999"
)

/*
	项目的 controllers 部分是负责提取指标并生成 collector 接口中的方法
	然后通过 registrer 注册
	最后生成
*/

func main() {
	// 打开数据库
	db, err := sql.Open(driver, dsn)
	if err != nil {
		log.Fatal("db open：", err)
	}
	defer db.Close()
	if err := db.Ping(); err != nil {
		log.Fatal("db ping：", err)
	}

	// 注册指标
	prometheus.Register(controllers.NewUpCollertor(db))
	prometheus.Register(controllers.NewThreadCoThreadsCollector(db))
	prometheus.Register(controllers.NewLongtimeLongtimeCollector(db))
	// 开启服务
	http.Handle("/metrics/", promhttp.Handler())

	// 启动 web 服务
	http.ListenAndServe(addr, nil)
}
