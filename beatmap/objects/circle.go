package objects

import (
	"github.com/Wieku/gosu-pp/beatmap/difficulty"
)

type Circle struct {
	*HitObject

	diff *difficulty.Difficulty
}

func NewCircle(data []string) *Circle {
	return &Circle{
		HitObject: commonParse(data),
	}
}

func (circle *Circle) GetType() Type {
	return CIRCLE
}
