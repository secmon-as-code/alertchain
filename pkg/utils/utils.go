package utils

func ToPtrSlice[T any](x []T) []*T {
	y := make([]*T, len(x))
	for i := range x {
		y[i] = &x[i]
	}
	return y
}
