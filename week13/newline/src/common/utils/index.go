package utils

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
)

func JSONToMap(str string) map[string]interface{} {
	var tempMap map[string]interface{}
	err := json.Unmarshal([]byte(str), &tempMap)
	if err != nil {
		panic(err)
	}

	return tempMap
}

func StructToMap(obj interface{}) map[string]interface{} {
	obj_v := reflect.ValueOf(obj)
	v := obj_v.Elem()
	typeOfType := v.Type()
	var data = make(map[string]interface{})
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		data[typeOfType.Field(i).Name] = field.Interface()
	}
	return data
}

func IsZero(f interface{}) bool {
	v := reflect.ValueOf(f)
	switch v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.String:
		str := v.String()
		if str == "" {
			return true
		}
		zero, error := strconv.ParseFloat(str, 10)
		if zero == 0 && error == nil {
			return true
		}
		boolean, error := strconv.ParseBool(str)
		return boolean == false && error == nil
	default:
		return false
	}
}

func IsNil(i interface{}) bool {
	vi := reflect.ValueOf(i)
	if vi.Kind() == reflect.Ptr {
		return vi.IsNil()
	}
	return false
}

func PaggingLoop(pageSize int32, lastOid string, startDate string, endDate string, cb func(int32, string, string, string) (int32, string)) {
	total, lastOid := cb(pageSize, lastOid, startDate, endDate)

	for total > 0 && lastOid != "" {
		total, lastOid = cb(pageSize, lastOid, startDate, endDate)
		fmt.Print("total, lastOid", total, lastOid)
	}
}

func ternary(flag bool, a interface{}, b interface{}) interface{} {
	if flag {
		return a
	} else {
		return b
	}
}

func And(a interface{}, b interface{}) interface{} {
	if !IsZero(a) {
		return b
	}
	return a
}

func Or(a interface{}, b interface{}) interface{} {
	if IsZero(a) {
		return b
	}
	return a
}

func Either(a interface{}, b interface{}) bool {
	if !IsZero(a) && !IsZero(b) {
		return true
	}
	return false
}

func Neither(a interface{}, b interface{}) bool {
	if IsZero(a) && IsZero(b) {
		return true
	}
	return false
}
