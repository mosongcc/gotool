package gslice

// Contains 数组中是否包含指定值
func Contains[T string | int | uint | int64 | uint64 | float64 | bool](slice []T, v T) bool {
	for _, val := range slice {
		if val == v {
			return true
		}
	}
	return false
}

// AppendUnique 数组追加值，已存在则不追加
func AppendUnique[T string | int | uint | int64 | uint64 | float64 | bool](slice []T, elems ...T) []T {
	for _, elem := range elems {
		if Contains(slice, elem) {
			continue
		}
		slice = append(slice, elem)
	}
	return slice
}
