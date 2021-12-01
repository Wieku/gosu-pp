package curves

import (
	"github.com/Wieku/gosu-pp/math/vector"
)

type Curve interface {
	PointAt(t float32) vector.Vector2f
	GetStartAngle() float32
	GetEndAngle() float32
	GetLength() float32
}
