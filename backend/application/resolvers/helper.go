package resolvers

import "stop-checker.com/application/schema"

func refList[T any](t []T) []*T {
	results := make([]*T, len(t))
	for i, result := range t {
		ref := result
		results[i] = &ref
	}
	return results
}

func ref[T any](t T) *T {
	return &t
}

func apply[T any, U any](t []T, f func(T) U) []U {
	results := make([]U, len(t))
	for i, result := range t {
		results[i] = f(result)
	}
	return results
}

func nullable[T any](t T, err error) *T {
	if err != nil {
		return nil
	}
	return &t
}

type Page[T any] struct {
	schema.PageInput
	data []T
}

func Paginate[T any](data []T, req schema.PageInput) *Page[T] {
	if req.Limit < 0 {
		req.Limit = len(data)
	}
	return &Page[T]{PageInput: req, data: data}
}

func (p *Page[T]) Info() *schema.PageInfo {
	return &schema.PageInfo{
		Cursor:    p.cursor(),
		Remaining: p.remaining(),
	}
}

func (p *Page[T]) Results() []T {
	if p.Skip >= len(p.data) {
		return []T{}
	}

	return p.data[p.Skip:p.cursor()]
}

func (p *Page[T]) remaining() int {
	return len(p.data) - p.cursor()
}

func (p *Page[T]) cursor() int {
	if p.Skip+p.Limit > len(p.data) {
		return len(p.data)
	}

	return p.Skip + p.Limit
}
