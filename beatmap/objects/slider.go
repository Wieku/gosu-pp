package objects

import (
	"github.com/Wieku/gosu-pp/beatmap/audio"
	"github.com/Wieku/gosu-pp/beatmap/difficulty"
	"github.com/Wieku/gosu-pp/beatmap/timing"
	"github.com/Wieku/gosu-pp/math/curves"
	"github.com/Wieku/gosu-pp/math/mutils"
	"github.com/Wieku/gosu-pp/math/vector"
	"math"
	"sort"
	"strconv"
	"strings"
)

type TickPoint struct {
	Time      float64
	IsReverse bool
}

type Slider struct {
	*HitObject

	multiCurve *curves.MultiCurve

	Timings *timing.Timings
	TPoint  timing.ControlPoint

	pixelLength float64
	RepeatCount int64

	ScorePoints []TickPoint

	diff *difficulty.Difficulty

	spanDuration float64
}

func NewSlider(data []string) *Slider {
	slider := &Slider{
		HitObject: commonParse(data),
	}

	slider.PositionDelegate = slider.PositionAt

	slider.pixelLength, _ = strconv.ParseFloat(data[7], 64)
	slider.RepeatCount, _ = strconv.ParseInt(data[6], 10, 64)

	list := strings.Split(data[5], "|")
	points := []vector.Vector2f{slider.StartPosRaw}

	for i := 1; i < len(list); i++ {
		list2 := strings.Split(list[i], ":")
		x, _ := strconv.ParseFloat(list2[0], 32)
		y, _ := strconv.ParseFloat(list2[1], 32)
		points = append(points, vector.NewVec2f(float32(x), float32(y)))
	}

	slider.multiCurve = curves.NewMultiCurveT(list[0], points, slider.pixelLength)

	slider.EndTime = slider.StartTime
	slider.EndPosRaw = slider.multiCurve.PointAt(1.0)

	baseSample := slider.sounds[0]

	slider.sounds = make([]audio.HitSound, slider.RepeatCount+1)

	for i := range slider.sounds {
		slider.sounds[i] = baseSample
	}

	if len(data) > 8 {
		subData := strings.Split(data[8], "|")
		for i, v := range subData {
			sound, _ := strconv.Atoi(v)
			slider.sounds[i] = audio.HitSound(sound)
		}
	}

	return slider
}

func (slider *Slider) PositionAt(time float64) vector.Vector2f {
	if slider.IsRetarded() {
		return slider.StartPosRaw
	}

	t1 := mutils.ClampF64(time, slider.StartTime, slider.EndTime)

	progress := (t1 - slider.StartTime) / slider.spanDuration

	progress = math.Mod(progress, 2)
	if progress >= 1 {
		progress = 2 - progress
	}

	return slider.multiCurve.PointAt(float32(progress))
}

func (slider *Slider) SetTiming(timings *timing.Timings) {
	slider.Timings = timings
	slider.TPoint = timings.GetPointAt(slider.StartTime)

	nanTimingPoint := math.IsNaN(slider.TPoint.GetRawBeatLength())

	velocity := slider.Timings.GetVelocity(slider.TPoint)

	cLength := float64(slider.multiCurve.GetLength())

	slider.spanDuration = cLength * 1000 / velocity

	slider.EndTime = slider.StartTime + cLength*1000*float64(slider.RepeatCount)/velocity

	minDistanceFromEnd := velocity * 0.01
	tickDistance := slider.Timings.GetTickDistance(slider.TPoint)

	if slider.multiCurve.GetLength() > 0 && tickDistance > slider.pixelLength {
		tickDistance = slider.pixelLength
	}

	for span := 0; span < int(slider.RepeatCount); span++ {
		spanStartTime := slider.StartTime + float64(span)*slider.spanDuration
		reversed := span%2 == 1

		// skip ticks if timingPoint has NaN beatLength
		for d := tickDistance; d <= cLength && !nanTimingPoint; d += tickDistance {
			if d >= cLength-minDistanceFromEnd {
				break
			}

			// Always generate ticks from the start of the path rather than the span to ensure that ticks in repeat spans are positioned identically to those in non-repeat spans
			timeProgress := d / cLength
			if reversed {
				timeProgress = 1 - timeProgress
			}

			slider.ScorePoints = append(slider.ScorePoints, TickPoint{
				Time: spanStartTime + timeProgress*slider.spanDuration,
			})
		}

		if span < int(slider.RepeatCount)-1 {
			slider.ScorePoints = append(slider.ScorePoints, TickPoint{
				Time:      spanStartTime + slider.spanDuration,
				IsReverse: true,
			})
		} else {
			slider.ScorePoints = append(slider.ScorePoints, TickPoint{
				Time: math.Max(slider.StartTime+(slider.EndTime-slider.StartTime)/2, slider.EndTime-36),
			})
		}
	}

	sort.Slice(slider.ScorePoints, func(i, j int) bool {
		return slider.ScorePoints[i].Time < slider.ScorePoints[j].Time
	})

	slider.EndPosRaw = slider.PositionAt(slider.EndTime)
}

func (slider *Slider) SetDifficulty(diff *difficulty.Difficulty) {
	slider.diff = diff
}

func (slider *Slider) IsRetarded() bool {
	return slider.StartTime == slider.EndTime
}

func (slider *Slider) GetType() Type {
	return SLIDER
}
