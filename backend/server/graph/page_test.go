package graph

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"stop-checker.com/server/graph/sdl"
)

func TestPage(t *testing.T) {
	data := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}

	t.Run("0 skip, 0 limit", func(t *testing.T) {
		page := Paginate(data, sdl.PageInput{
			Skip: 0, Limit: 0,
		})

		info := page.Info()
		assert.Equal(t, 10, info.Remaining)
		assert.Equal(t, 0, info.Cursor)

		results := page.Results()
		assert.Equal(t, 0, len(results))
	})

	t.Run("0 skip, 6 limit", func(t *testing.T) {
		page := Paginate(data, sdl.PageInput{
			Skip: 0, Limit: 6,
		})

		info := page.Info()
		assert.Equal(t, 4, info.Remaining)
		assert.Equal(t, 6, info.Cursor)

		results := page.Results()
		assert.Equal(t, 6, len(results))
		assert.Equal(t, 0, results[0])
		assert.Equal(t, 5, results[5])
	})

	t.Run("4 skip, 6 limit", func(t *testing.T) {
		page := Paginate(data, sdl.PageInput{
			Skip: 4, Limit: 6,
		})

		info := page.Info()
		assert.Equal(t, 0, info.Remaining)
		assert.Equal(t, 10, info.Cursor)

		results := page.Results()
		assert.Equal(t, 6, len(results))
	})

	t.Run("10 skip, 1 limit", func(t *testing.T) {
		page := Paginate(data, sdl.PageInput{
			Skip: 10, Limit: 1,
		})

		info := page.Info()
		assert.Equal(t, 0, info.Remaining)
		assert.Equal(t, 10, info.Cursor)

		results := page.Results()
		assert.Equal(t, 0, len(results))
	})

	t.Run("5 skip, 10 limit", func(t *testing.T) {
		page := Paginate(data, sdl.PageInput{
			Skip: 5, Limit: 10,
		})

		info := page.Info()
		assert.Equal(t, 0, info.Remaining)
		assert.Equal(t, 10, info.Cursor)

		results := page.Results()
		assert.Equal(t, 5, len(results))
		assert.Equal(t, 5, results[0])
		assert.Equal(t, 9, results[4])
	})

	t.Run("5 skip, -1 limit", func(t *testing.T) {
		page := Paginate(data, sdl.PageInput{
			Skip: 5, Limit: -1,
		})

		info := page.Info()
		assert.Equal(t, 0, info.Remaining)
		assert.Equal(t, 10, info.Cursor)

		results := page.Results()
		assert.Equal(t, 5, len(results))
		assert.Equal(t, 5, results[0])
		assert.Equal(t, 9, results[4])
	})
}
