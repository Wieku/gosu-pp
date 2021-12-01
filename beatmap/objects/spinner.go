package objects

import (
	"github.com/Wieku/gosu-pp/beatmap/timing"
	"strconv"
)

type Spinner struct {
	*HitObject

	Timings *timing.Timings
}

func NewSpinner(data []string) *Spinner {
	spinner := &Spinner{
		HitObject: commonParse(data),
	}

	spinner.EndTime, _ = strconv.ParseFloat(data[5], 64)

	return spinner
}

func (spinner *Spinner) SetTiming(timings *timing.Timings) {
	spinner.Timings = timings
}

func (spinner *Spinner) GetType() Type {
	return SPINNER
}
