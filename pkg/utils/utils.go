package utils

func ToPtrSlice[T any](x []T) []*T {
	y := make([]*T, len(x))
	for i, v := range x {
		y[i] = &v
	}
	return y
}
