package utils

import (
	"reflect"
	"strconv"
	"time"
	"errors"
)

func MapToStructByTagSql(data map[string]string, obj interface{}) (err error) {
	objValues := reflect.ValueOf(obj).Elem()
	for i := 0; i < objValues.NumField(); i++ {
		value := data[objValues.Type().Field(i).Tag.Get("sql")]
		/* 获取sql字段查询出来的值 */
		name := objValues.Type().Field(i).Name
		/* 获取数据模型字段名称 */
		structFieldType := objValues.Field(i).Type()
		/* 获取数据模型字段类型 */
		val := reflect.ValueOf(value)
		/* 获取sql查询出来的结构体的值类型 */
		if structFieldType != val.Type() {
			/* 数据库存储类型和定义的结构体字段类型可能不同，所以我们要进行类型转换之后再赋值 */
			val, err = TypeConversion(value, structFieldType.Name())
			if err != nil {
				return
			}
		}
		objValues.FieldByName(name).Set(val)
		/* 填充数据模型 */
	}
	return
}

//类型转换
func TypeConversion(value string, ntype string) (reflect.Value, error) {
	if ntype == "string" {
		return reflect.ValueOf(value), nil
	} else if ntype == "time.Time" {
		t, err := time.ParseInLocation("2006-01-02 15:04:05", value, time.Local)
		return reflect.ValueOf(t), err
	} else if ntype == "Time" {
		t, err := time.ParseInLocation("2006-01-02 15:04:05", value, time.Local)
		return reflect.ValueOf(t), err
	} else if ntype == "int" {
		i, err := strconv.Atoi(value)
		return reflect.ValueOf(i), err
	} else if ntype == "int8" {
		i, err := strconv.ParseInt(value, 10, 64)
		return reflect.ValueOf(int8(i)), err
	} else if ntype == "int32" {
		i, err := strconv.ParseInt(value, 10, 64)
		return reflect.ValueOf(int64(i)), err
	} else if ntype == "int64" {
		i, err := strconv.ParseInt(value, 10, 64)
		return reflect.ValueOf(i), err
	} else if ntype == "float32" {
		i, err := strconv.ParseFloat(value, 64)
		return reflect.ValueOf(float32(i)), err
	} else if ntype == "float64" {
		i, err := strconv.ParseFloat(value, 64)
		return reflect.ValueOf(i), err
	}

	//else if .......增加其他一些类型的转换

	return reflect.ValueOf(value), errors.New("未知的类型：" + ntype)
}