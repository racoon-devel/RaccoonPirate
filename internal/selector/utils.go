package selector

func getValue[T int | uint | int32 | uint32 | int64 | uint64](val *T) T {
	if val == nil {
		return T(0)
	}

	return *val
}

func getString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

type ordered interface {
	int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64 | float32 | float64
}

func findMax[T any, C ordered](arr []T, f func(elem T) C) (T, C, int) {
	if len(arr) == 0 {
		panic("array must be non empty")
	}
	max := arr[0]
	index := 0
	for i := range arr {
		if f(arr[i]) > f(max) {
			max = arr[i]
			index = i
		}
	}

	return max, f(max), index
}
