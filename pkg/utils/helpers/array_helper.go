package helpers

func Filter[T any](s []T, cond func(t T) bool) []T {
	res := []T{}
	for _, v := range s {
		if cond(v) {
			res = append(res, v)
		}
	}
	return res
}

func Map[T, U any](s []T, f func(t T) U) []U {
	res := make([]U, len(s))
	for i, v := range s {
		res[i] = f(v)
	}
	return res
}
