package pointerutil

func MakePointer[T any](data T) *T {
	value := data
	return &value
}

func ResolveValue[T any](data *T) T {
	var defaultValue T
	if data == nil {
		return defaultValue
	}
	return *data
}

func ResolveNestedArray[T1 any, T2 any](slice *[]T1, converterFn func(T1) T2) (result []T2) {
	if slice == nil {
		return
	}
	result = make([]T2, 0, len(*slice))
	for _, item := range *slice {
		result = append(result, converterFn(item))
	}
	return result
}
