package graph

import "stop-checker.com/server/graph/sdl"

type Page[T any] struct {
	sdl.PageInput
	data []T
}

func Paginate[T any](data []T, req sdl.PageInput) *Page[T] {
	if req.Limit < 0 {
		req.Limit = len(data)
	}
	return &Page[T]{PageInput: req, data: data}
}

func (p *Page[T]) Info() *sdl.PageInfo {
	return &sdl.PageInfo{
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
