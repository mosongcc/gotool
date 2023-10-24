package gmap

// Keys 读取 map key 数组
func Keys[K string | int | int64, V any](m map[K]V) (keys []K) {
	for k, _ := range m {
		keys = append(keys, k)
	}
	return
}
