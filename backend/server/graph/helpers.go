package graph

func ref[T any](t []T) []*T {
	results := make([]*T, len(t))
	for i, result := range t {
		ref := result
		results[i] = &ref
	}
	return results
}
func apply[T any, U any](t []T, f func(T) U) []U {
	results := make([]U, len(t))
	for i, result := range t {
		results[i] = f(result)
	}
	return results
}
