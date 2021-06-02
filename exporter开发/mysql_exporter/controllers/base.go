package controllers

import (
	"database/sql"
)

type mysqlSelcetor struct {
	mysqlconnection *sql.DB
}

// 构建函数
func newMysqlSelcector(sqldb *sql.DB) *mysqlSelcetor {
	return &mysqlSelcetor{
		mysqlconnection: sqldb,
	}
}

// status 查询的方法
func (selector *mysqlSelcetor) showStatus(value string) float64 {
	sqlrow := selector.mysqlconnection.QueryRow("SHOW GLOBAL status WHERE variable_name=?;", value)
	var (
		vname  string
		vvalue float64
	)
	sqlrow.Scan(&vname, &vvalue)
	return vvalue
}

// variable 查询的方法，返回 float64
func (selector *mysqlSelcetor) showVariable(value string) float64 {
	sqlrow := selector.mysqlconnection.QueryRow("show global variables where variable_name=?", value)
	var (
		vname  string
		vvalue float64
	)
	sqlrow.Scan(&vname, &vvalue)
	return vvalue
}
