package response

import "reflect"

// NormalizeArrayIfNeeded return empty array if v is a nil slice
func NormalizeArrayIfNeeded(v interface{}) interface{} {
	val := reflect.ValueOf(v)
	if (val.Kind() == reflect.Array || val.Kind() == reflect.Slice) && val.Len() == 0 {
		return make([]int, 0)
	}
	return v
}
