package utils

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"reflect"
	"strconv"
)

func NewMysqlConn() (db *sql.DB,err error) {
	dsn := "root:root@tcp(127.0.0.1:3306)/kapi"
	return sql.Open("mysql", dsn)
}

func SQLToMap(rows *sql.Rows) map[string]string {
	cols,_ :=  rows.Columns()
	scanArgs := make([]interface{}, len(cols))
	values := make([]interface{}, len(cols))

	for idx := range values {
		scanArgs[idx] = &values[idx]
	}

	record := make(map[string]string)

	for rows.Next() {
		rows.Scan(scanArgs...)
		for idx,item := range values {
			if item != nil {
				if reflect.TypeOf(item).Kind() != reflect.Int64 {
					record[cols[idx]] = string(item.([]byte))
				}else{
					record[cols[idx]] = strconv.Itoa(int(item.(int64)))
				}
			}
		}
	}
	return record
}