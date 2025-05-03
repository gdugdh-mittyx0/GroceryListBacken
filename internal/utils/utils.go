package utils

import (
	"strconv"
)

func SerializeBool(b bool) string {
	if b {
		return "true"
	}
	return "false"
}

func SerializeUintToMapIfNotempty(m *map[string]string, key string, val uint64) {
	if val != 0 {
		(*m)[key] = strconv.FormatUint(val, 10)
	}
}

func GetInterfacePtr[T any](v interface{}) *T {
	if val, ok := v.(T); ok {
		return &val
	}
	if val, ok := v.(*T); ok {
		return val
	}
	return nil
}

func GetInterfaceValue[T any](v interface{}) T {
	if val, ok := v.(T); ok {
		return val
	}
	val := new(T)
	return *val
}

func ConvertToInterfaceSlice[T any](data []T) []interface{} {
	result := make([]interface{}, len(data))
	for i, v := range data {
		result[i] = v
	}
	return result
}
