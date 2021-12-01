package beatmap

import (
	"github.com/Wieku/gosu-pp/beatmap/difficulty"
	"github.com/Wieku/gosu-pp/beatmap/objects"
	"github.com/Wieku/gosu-pp/beatmap/timing"
	"math"
	"strconv"
	"strings"
)

type BeatMap struct {
	FileVersion int

	Artist        string
	ArtistUnicode string

	Title        string
	TitleUnicode string

	Version string

	Creator string

	Source string

	Tags string

	Mode int

	SliderMultiplier float64
	StackLeniency    float64

	Difficulty *difficulty.Difficulty

	Audio string
	Bg    string

	SetID int64
	MapID int64

	PreviewTime int64

	Length   int
	Circles  int
	Sliders  int
	Spinners int

	MinBPM float64
	MaxBPM float64

	Timings    *timing.Timings
	HitObjects []objects.IHitObject
	Pauses     []*Pause

	arSpecified bool
}

func NewBeatMap() *BeatMap {
	return &BeatMap{
		Timings:       timing.NewTimings(),
		StackLeniency: 0.7,
		Difficulty:    difficulty.NewDifficulty(5, 5, 5, 5),

		MinBPM: math.Inf(0),
		MaxBPM: 0,
	}
}

func (beatMap *BeatMap) ParsePoint(point string) {
	line := strings.Split(point, ",")
	pointTime, _ := strconv.ParseInt(line[0], 10, 64)
	bpm, _ := strconv.ParseFloat(line[1], 64)

	if !math.IsNaN(bpm) && bpm >= 0 {
		rBPM := 60000 / bpm
		beatMap.MinBPM = math.Min(beatMap.MinBPM, rBPM)
		beatMap.MaxBPM = math.Max(beatMap.MaxBPM, rBPM)
	}

	signature := 4
	sampleSet := 0 //beatMap.Timings.LastSet
	sampleIndex := 1
	sampleVolume := 1.0
	inherited := false
	kiai := false
	omitFirstBarLine := false

	if len(line) > 2 {
		signature, _ = strconv.Atoi(line[2])
		if signature == 0 {
			signature = 4
		}
	}

	if len(line) > 3 {
		sampleSet, _ = strconv.Atoi(line[3])
	}

	if len(line) > 4 {
		sampleIndex, _ = strconv.Atoi(line[4])
	}

	if len(line) > 5 {
		sV, _ := strconv.Atoi(line[5])
		sampleVolume = float64(sV) / 100
	}

	if len(line) > 6 {
		inh, _ := strconv.Atoi(line[6])
		inherited = inh == 0
	}

	if len(line) > 7 {
		ki, _ := strconv.Atoi(line[7])
		kiai = (ki & 1) > 0
		omitFirstBarLine = (ki & 8) > 0
	}

	beatMap.Timings.AddPoint(float64(pointTime), bpm, sampleSet, sampleIndex, sampleVolume, signature, inherited, kiai, omitFirstBarLine)
}

func (beatMap *BeatMap) FinalizePoints() {
	beatMap.Timings.FinalizePoints()
}
