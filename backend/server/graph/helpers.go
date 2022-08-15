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

func nullable[T any](t T, notFound error) (*T, error) {
	if notFound != nil {
		return nil, nil
	}
	return &t, nil
}

func nullableRef[T *any](t T, err error) (T, error) {
	if err != nil {
		return nil, nil
	}
	return t, nil
}
