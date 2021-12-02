package osu

import (
	"github.com/Wieku/gosu-pp/beatmap/difficulty"
	"github.com/Wieku/gosu-pp/beatmap/objects"
	"github.com/Wieku/gosu-pp/performance/osu/preprocessing"
	"github.com/Wieku/gosu-pp/performance/osu/skills"
	"math"
)

const (
	// StarScalingFactor is a global stars multiplier
	StarScalingFactor float64 = 0.0675
)

type Attributes struct {
	// Total Star rating, visible on osu!'s beatmap page
	Total float64

	// Aim stars, needed for Performance Points (aka PP) calculations
	Aim float64

	// SliderFactor is a ratio of Aim calculated without sliders to Aim with them
	SliderFactor float64

	// Speed stars, needed for Performance Points (aka PP) calculations
	Speed float64

	// Flashlight stars, needed for Performance Points (aka PP) calculations
	Flashlight float64

	ObjectCount int
	Circles     int
	Sliders     int
	Spinners    int
	MaxCombo    int
}

// StrainPeaks contains peaks of Aim, Speed and Flashlight skills, as well as peaks passed through star rating formula
type StrainPeaks struct {
	// Aim peaks
	Aim []float64

	// Speed peaks
	Speed []float64

	// Flashlight peaks
	Flashlight []float64

	// Total contains aim, speed and flashlight peaks passed through star rating formula
	Total []float64
}

// getStarsFromRawValues converts raw skill values to Attributes
func getStarsFromRawValues(rawAim, rawAimNoSliders, rawSpeed, rawFlashlight float64, diff *difficulty.Difficulty, attr Attributes) Attributes {
	aimRating := math.Sqrt(rawAim) * StarScalingFactor
	aimRatingNoSliders := math.Sqrt(rawAimNoSliders) * StarScalingFactor
	speedRating := math.Sqrt(rawSpeed) * StarScalingFactor
	flashlightVal := math.Sqrt(rawFlashlight) * StarScalingFactor

	sliderFactor := 1.0
	if aimRating > 0.00001 {
		sliderFactor = aimRatingNoSliders / aimRating
	}

	var total float64

	if diff.CheckModActive(difficulty.Relax) {
		speedRating = 0.0
	}

	baseAimPerformance := ppBase(aimRating)
	baseSpeedPerformance := ppBase(speedRating)
	baseFlashlightPerformance := 0.0

	if diff.CheckModActive(difficulty.Flashlight) {
		baseFlashlightPerformance = math.Pow(flashlightVal, 2.0) * 25.0
	}

	basePerformance := math.Pow(
		math.Pow(baseAimPerformance, 1.1)+
			math.Pow(baseSpeedPerformance, 1.1)+
			math.Pow(baseFlashlightPerformance, 1.1),
		1.0/1.1,
	)

	if basePerformance > 0.00001 {
		total = math.Cbrt(1.12) * 0.027 * (math.Cbrt(100000/math.Pow(2, 1/1.1)*basePerformance) + 4)
	}

	attr.Total = total
	attr.Aim = aimRating
	attr.SliderFactor = sliderFactor
	attr.Speed = speedRating
	attr.Flashlight = flashlightVal

	return attr
}

// Retrieves skill values and converts to Attributes
func getStars(aim *skills.AimSkill, aimNoSliders *skills.AimSkill, speed *skills.SpeedSkill, flashlight *skills.Flashlight, diff *difficulty.Difficulty, attr Attributes) Attributes {
	rawFlashlight := 0.0
	if flashlight != nil {
		rawFlashlight = flashlight.DifficultyValue()
	}

	return getStarsFromRawValues(
		aim.DifficultyValue(),
		aimNoSliders.DifficultyValue(),
		speed.DifficultyValue(),
		rawFlashlight,
		diff,
		attr,
	)
}

// Retrieves peaks from skills
func getPeaks(aim *skills.AimSkill, speed *skills.SpeedSkill, flashlight *skills.Flashlight, diff *difficulty.Difficulty) StrainPeaks {
	peaks := StrainPeaks{
		Aim:   aim.GetCurrentStrainPeaks(),
		Speed: speed.GetCurrentStrainPeaks(),
	}

	if flashlight != nil {
		peaks.Flashlight = flashlight.GetCurrentStrainPeaks()
	}

	peaks.Total = make([]float64, len(peaks.Aim))

	for i := 0; i < len(peaks.Aim); i++ {
		flVal := 0.0
		if flashlight != nil {
			flVal = peaks.Flashlight[i]
		}

		stars := getStarsFromRawValues(peaks.Aim[i], peaks.Aim[i], peaks.Speed[i], flVal, diff, Attributes{})
		peaks.Total[i] = stars.Total
	}

	return peaks
}

func addObjectToAttribs(o objects.IHitObject, attr *Attributes) {
	if s, ok := o.(*objects.Slider); ok {
		attr.Sliders++
		attr.MaxCombo += len(s.ScorePoints)
	} else if _, ok := o.(*objects.Circle); ok {
		attr.Circles++
	} else if _, ok := o.(*objects.Spinner); ok {
		attr.Spinners++
	}

	attr.MaxCombo++
	attr.ObjectCount++
}

func calculateSingle(objects []objects.IHitObject, diff *difficulty.Difficulty) (Attributes, *skills.AimSkill, *skills.SpeedSkill, *skills.Flashlight) {
	diffObjects := preprocessing.CreateDifficultyObjects(objects, diff)

	aimSkill := skills.NewAimSkill(diff, true)
	aimNoSlidersSkill := skills.NewAimSkill(diff, false)
	speedSkill := skills.NewSpeedSkill(diff)

	var flashlightSkill *skills.Flashlight
	if diff.CheckModActive(difficulty.Flashlight) {
		flashlightSkill = skills.NewFlashlightSkill(diff)
	}

	attr := Attributes{}

	addObjectToAttribs(objects[0], &attr)

	for i, o := range diffObjects {
		addObjectToAttribs(objects[i+1], &attr)

		aimSkill.Process(o)
		aimNoSlidersSkill.Process(o)
		speedSkill.Process(o)

		if flashlightSkill != nil {
			flashlightSkill.Process(o)
		}
	}

	stars := getStars(aimSkill, aimNoSlidersSkill, speedSkill, flashlightSkill, diff, attr)

	return stars, aimSkill, speedSkill, flashlightSkill
}

// CalculateSingle calculates the final difficulty attributes of a map
func CalculateSingle(objects []objects.IHitObject, diff *difficulty.Difficulty) Attributes {
	stars, _, _, _ := calculateSingle(objects, diff)
	return stars
}

// CalculateSingleWithPeaks calculates the final difficulty attributes, and strain peaks of a map
func CalculateSingleWithPeaks(objects []objects.IHitObject, diff *difficulty.Difficulty) (Attributes, StrainPeaks) {
	stars, aimSkill, speedSkill, flashlightSkill := calculateSingle(objects, diff)
	return stars, getPeaks(aimSkill, speedSkill, flashlightSkill, diff)
}

func calculateStep(objects []objects.IHitObject, diff *difficulty.Difficulty) ([]Attributes, *skills.AimSkill, *skills.SpeedSkill, *skills.Flashlight) {
	diffObjects := preprocessing.CreateDifficultyObjects(objects, diff)

	aimSkill := skills.NewAimSkill(diff, true)
	aimNoSlidersSkill := skills.NewAimSkill(diff, false)
	speedSkill := skills.NewSpeedSkill(diff)

	var flashlightSkill *skills.Flashlight
	if diff.CheckModActive(difficulty.Flashlight) {
		flashlightSkill = skills.NewFlashlightSkill(diff)
	}

	stars := make([]Attributes, 1, len(objects))

	addObjectToAttribs(objects[0], &stars[0])

	for i, o := range diffObjects {
		attr := stars[i]
		addObjectToAttribs(objects[i+1], &attr)

		aimSkill.Process(o)
		aimNoSlidersSkill.Process(o)
		speedSkill.Process(o)

		if flashlightSkill != nil {
			flashlightSkill.Process(o)
		}

		stars = append(stars, getStars(aimSkill, aimNoSlidersSkill, speedSkill, flashlightSkill, diff, attr))
	}

	return stars, aimSkill, speedSkill, flashlightSkill
}

// CalculateStep calculates successive star ratings for every part of a beatmap
func CalculateStep(objects []objects.IHitObject, diff *difficulty.Difficulty) []Attributes {
	stars, _, _, _ := calculateStep(objects, diff)
	return stars
}

// CalculateStepWithPeaks calculates strain peaks and successive star ratings for every part of a beatmap
func CalculateStepWithPeaks(objects []objects.IHitObject, diff *difficulty.Difficulty) ([]Attributes, StrainPeaks) {
	stars, aimSkill, speedSkill, flashlightSkill := calculateStep(objects, diff)
	return stars, getPeaks(aimSkill, speedSkill, flashlightSkill, diff)
}

// CalculateStrainPeaks calculates difficulty strain peaks of a beatmap
func CalculateStrainPeaks(objects []objects.IHitObject, diff *difficulty.Difficulty) StrainPeaks {
	diffObjects := preprocessing.CreateDifficultyObjects(objects, diff)

	aimSkill := skills.NewAimSkill(diff, true)
	speedSkill := skills.NewSpeedSkill(diff)

	var flashlightSkill *skills.Flashlight
	if diff.CheckModActive(difficulty.Flashlight) {
		flashlightSkill = skills.NewFlashlightSkill(diff)
	}

	for _, o := range diffObjects {
		aimSkill.Process(o)
		speedSkill.Process(o)

		if flashlightSkill != nil {
			flashlightSkill.Process(o)
		}
	}

	return getPeaks(aimSkill, speedSkill, flashlightSkill, diff)
}
