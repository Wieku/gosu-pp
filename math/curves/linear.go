package curves

import (
	"github.com/Wieku/gosu-pp/math/vector"
)

type Linear struct {
	Point1, Point2 vector.Vector2f
}

func NewLinear(pt1, pt2 vector.Vector2f) Linear {
	return Linear{pt1, pt2}
}

func (ln Linear) PointAt(t float32) vector.Vector2f {
	return ln.Point1.Lerp(ln.Point2, t)
}

func (ln Linear) GetStartAngle() float32 {
	return ln.Point1.AngleRV(ln.Point2)
}

func (ln Linear) GetEndAngle() float32 {
	return ln.Point2.AngleRV(ln.Point1)
}

func (ln Linear) GetLength() float32 {
	return ln.Point1.Dst(ln.Point2)
}
