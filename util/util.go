package util

import "reflect"

func IsDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func IsAlphabet(ch byte) bool {
	return ('a' <= ch && ch <= 'z') || ('A' <= ch && ch <= 'Z')
}

func IsInt(num any) bool {
	switch reflect.ValueOf(num).Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return true
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return true
	default:
		return false
	}
}

func IsFloat(num any) bool {
	switch reflect.ValueOf(num).Kind() {
	case reflect.Float32, reflect.Float64:
		return true
	default:
		return false
	}
}

func IsString(num any) bool {
	switch reflect.ValueOf(num).Kind() {
	case reflect.String:
		return true
	default:
		return false
	}
}

func ConvertToInt(num any) int64 {
	val := reflect.ValueOf(num)
	switch val.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return val.Int()
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return int64(val.Uint())
	default:
		return 0
	}
}

func ConvertToFloat(num any) float64 {
	val := reflect.ValueOf(num)
	switch val.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return float64(val.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return float64(val.Uint())
	case reflect.Float32, reflect.Float64:
		return val.Float()
	default:
		return 0
	}
}

func ConvertToBool(val reflect.Value) bool {
	switch val.Kind() {
	case reflect.Bool:
		return val.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return val.Int() > 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return val.Uint() > 0
	case reflect.Float32, reflect.Float64:
		return val.Float() > 0
	case reflect.String:
		return val.String() != ""
	default:
		return false
	}
}

func ConvertToString(num any) string {
	val := reflect.ValueOf(num)
	switch val.Kind() {
	case reflect.String:
		return val.String()
	default:
		return ""
	}
}
