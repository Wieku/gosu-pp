package test

import (
	"github.com/wieku/gosu-pp/beatmap"
	"github.com/wieku/gosu-pp/performance/osu"
	"os"
	"testing"
)

func TestOsuSR(t *testing.T) {
	osuFile, err := os.Open("inabakumori - Lost Umbrella (Ryuusei Aika) [________].osu")
	if err != nil {
		t.Fatal(err)
	}

	beatMap, err1 := beatmap.ParseBeatMapFile(osuFile)
	if err1 != nil {
		t.Fatal(err)
	}

	stars := osu.CalculateSingle(beatMap.HitObjects, beatMap.Diff)

	AssertFloat(4.633597654263426, stars.Aim, 0.0001, t)
	AssertFloat(8.234703277960627, stars.Total, 0.0001, t)
	AssertFloat(2.812916095444383, stars.Speed, 0.0001, t)
}

func BenchmarkOsuSRSingle(b *testing.B) {
	osuFile, err := os.Open("Avenged Sevenfold - Save Me (Drummer) [Tragedy].osu")
	if err != nil {
		b.Fatal(err)
	}

	beatMap, err1 := beatmap.ParseBeatMapFile(osuFile)
	if err1 != nil {
		b.Fatal(err)
	}

	for n := 0; n < b.N; n++ {
		osu.CalculateSingle(beatMap.HitObjects, beatMap.Diff)
	}
}

func BenchmarkOsuSRStep(b *testing.B) {
	osuFile, err := os.Open("Avenged Sevenfold - Save Me (Drummer) [Tragedy].osu")
	if err != nil {
		b.Fatal(err)
	}

	beatMap, err1 := beatmap.ParseBeatMapFile(osuFile)
	if err1 != nil {
		b.Fatal(err)
	}

	for n := 0; n < b.N; n++ {
		osu.CalculateStep(beatMap.HitObjects, beatMap.Diff)
	}
}
