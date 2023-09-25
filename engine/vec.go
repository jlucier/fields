package engine

import (
	"golang.org/x/exp/constraints"
	"math"
)

type GVec2[T constraints.Float] struct {
	X T
	Y T
}

func (self *GVec2[T]) Angle() float64 {
	return math.Atan2(float64(self.Y), float64(self.X))
}

// ops

func (self *GVec2[T]) Add(o GVec2[T]) GVec2[T] {
	return GVec2[T]{
		self.X + o.X,
		self.Y + o.Y,
	}
}

func (self *GVec2[T]) Sub(o GVec2[T]) GVec2[T] {
	return GVec2[T]{
		self.X - o.X,
		self.Y - o.Y,
	}
}

func (self *GVec2[T]) Mul(v T) GVec2[T] {
	return GVec2[T]{
		self.X * v,
		self.Y * v,
	}
}

func (self *GVec2[T]) Div(v T) GVec2[T] {
	return self.Mul(1 / v)
}

func (self *GVec2[T]) Norm() GVec2[T] {
	return self.Div(T(self.Mag()))
}

func (self *GVec2[T]) Mag() float64 {
	return math.Sqrt(math.Pow(float64(self.X), 2) + math.Pow(float64(self.Y), 2))
}

func (self *GVec2[T]) Dot(o GVec2[T]) T {
	return self.X*o.X + self.Y + o.Y
}

func Vec2FromAngle(a float64) Vec2 {
	v := Vec2{1, math.Tan(a)}
	return v.Norm()
}

type Vec2 = GVec2[float64]
