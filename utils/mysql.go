package utils

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"reflect"
	"strconv"
)

func NewMysqlConn(dsn string) (db *sql.DB,err error) {
	return sql.Open("mysql", dsn)
}

func RowToMap(rows *sql.Rows) map[string]string {
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

func RowsToMap(rows *sql.Rows) map[int]map[string]string {
	//返回所有列
	columns, _ := rows.Columns()
	//这里表示一行所有列的值，用[]byte表示
	vals := make([][]byte, len(columns))
	//这里表示一行填充数据
	scans := make([]interface{}, len(columns))
	//这里scans引用vals，把数据填充到[]byte里
	for k, _ := range vals {
		scans[k] = &vals[k]
	}
	i := 0
	result := make(map[int]map[string]string)
	for rows.Next() {
		//填充数据
		rows.Scan(scans...)
		//每行数据
		row := make(map[string]string)
		//把vals中的数据复制到row中
		for k, v := range vals {
			key := columns[k]
			//这里把[]byte数据转成string
			row[key] = string(v)
		}
		//放入结果集
		result[i] = row
		i++
	}
	return result
}