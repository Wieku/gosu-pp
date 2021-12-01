package objects

import (
	"github.com/Wieku/gosu-pp/beatmap/audio"
	"github.com/Wieku/gosu-pp/math/vector"
	"strconv"
)

func commonParse(data []string) *HitObject {
	x, _ := strconv.ParseFloat(data[0], 32)
	y, _ := strconv.ParseFloat(data[1], 32)
	time, _ := strconv.ParseFloat(data[2], 64)
	objType, _ := strconv.Atoi(data[3])

	startPos := vector.NewVec2f(float32(x), float32(y))

	sound, _ := strconv.Atoi(data[4])

	hitObject := &HitObject{
		StartPosRaw: startPos,
		EndPosRaw:   startPos,
		StartTime:   time,
		EndTime:     time,
		HitObjectID: -1,
		NewCombo:    (Type(objType) & NEWCOMBO) == NEWCOMBO,
		ColorOffset: (objType >> 4) & 7,
		sounds:      []audio.HitSound{audio.HitSound(sound)},
	}

	return hitObject
}
