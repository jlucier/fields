package engine

import (
	"github.com/veandco/go-sdl2/sdl"
)

func DrawCircle(renderer *sdl.Renderer, x0 int32, y0 int32, r int32) {
	x, y, dx, dy := r-1, int32(0), int32(1), int32(1)
	err := dx - (r * 2)

	for x > y {
		renderer.DrawPoint(x0+x, y0+y)
		renderer.DrawPoint(x0+y, y0+x)
		renderer.DrawPoint(x0-y, y0+x)
		renderer.DrawPoint(x0-x, y0+y)
		renderer.DrawPoint(x0-x, y0-y)
		renderer.DrawPoint(x0-y, y0-x)
		renderer.DrawPoint(x0+y, y0-x)
		renderer.DrawPoint(x0+x, y0-y)

		if err <= 0 {
			y++
			err += dy
			dy += 2
		}
		if err > 0 {
			x--
			dx += 2
			err += dx - (r * 2)
		}
	}

}
