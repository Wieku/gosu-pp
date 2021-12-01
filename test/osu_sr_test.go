package test

import (
	"github.com/Wieku/gosu-pp/beatmap"
	"github.com/Wieku/gosu-pp/beatmap/difficulty"
	"github.com/Wieku/gosu-pp/performance/osu"
	"os"
	"testing"
)

//TODO: Replace test values with osu-tools generated ones

func TestOsuSR(t *testing.T) {
	osuFile, err := os.Open("inabakumori - Lost Umbrella (Ryuusei Aika) [________].osu")
	if err != nil {
		t.Fatal(err)
	}

	beatMap, err1 := beatmap.ParseFromReader(osuFile)
	if err1 != nil {
		t.Fatal(err)
	}

	stars := osu.CalculateSingle(beatMap.HitObjects, beatMap.Difficulty)

	AssertFloat("Aim", 4.633597654263426, stars.Aim, 0.0001, t)
	AssertFloat("Speed", 2.812916095444383, stars.Speed, 0.0001, t)
	AssertFloat("Total", 8.234703277960627, stars.Total, 0.0001, t)
}

func TestOsuPP(t *testing.T) {
	osuFile, err := os.Open("Avenged Sevenfold - Save Me (Drummer) [Tragedy].osu")
	if err != nil {
		t.Fatal(err)
	}

	beatMap, err1 := beatmap.ParseFromReader(osuFile)
	if err1 != nil {
		t.Fatal(err)
	}

	stars := osu.CalculateSingle(beatMap.HitObjects, beatMap.Difficulty)

	pp := &osu.PPv2{}
	pp.PPv2x(stars, -1, -1, 0, 0, 0, beatMap.Difficulty)

	AssertFloat("Aim", 170.74687153986574, pp.Results.Aim, 0.0001, t)
	AssertFloat("Speed", 173.08993579731623, pp.Results.Speed, 0.0001, t)
	AssertFloat("Acc", 130.85913858310036, pp.Results.Acc, 0.0001, t)
	AssertFloat("Total", 481.4974251130584, pp.Results.Total, 0.0001, t)
}

func TestOsuPPDT(t *testing.T) {
	osuFile, err := os.Open("Avenged Sevenfold - Save Me (Drummer) [Tragedy].osu")
	if err != nil {
		t.Fatal(err)
	}

	beatMap, err1 := beatmap.ParseFromReader(osuFile)
	if err1 != nil {
		t.Fatal(err)
	}

	beatMap.Difficulty.SetMods(difficulty.DoubleTime)

	stars := osu.CalculateSingle(beatMap.HitObjects, beatMap.Difficulty)

	pp := &osu.PPv2{}
	pp.PPv2x(stars, stars.MaxCombo, stars.ObjectCount, 0, 0, 0, beatMap.Difficulty)

	AssertFloat("Aim", 569.94426438669574964, pp.Results.Aim, 0.0001, t)
	AssertFloat("Speed", 875.1324756257907, pp.Results.Speed, 0.0001, t)
	AssertFloat("Acc", 246.77051962773746, pp.Results.Acc, 0.0001, t)
	AssertFloat("Total", 1733.49497453401318126, pp.Results.Total, 0.0001, t)
}

func TestOsuPPHR(t *testing.T) {
	osuFile, err := os.Open("Avenged Sevenfold - Save Me (Drummer) [Tragedy].osu")
	if err != nil {
		t.Fatal(err)
	}

	beatMap, err1 := beatmap.ParseFromReader(osuFile)
	if err1 != nil {
		t.Fatal(err)
	}

	beatMap.Difficulty.SetMods(difficulty.HardRock)

	stars := osu.CalculateSingle(beatMap.HitObjects, beatMap.Difficulty)

	pp := &osu.PPv2{}
	pp.PPv2x(stars, stars.MaxCombo, stars.ObjectCount, 0, 0, 0, beatMap.Difficulty)

	AssertFloat("Aim", 219.08610212044914, pp.Results.Aim, 0.0001, t)
	AssertFloat("Speed", 195.3061931864951, pp.Results.Speed, 0.0001, t)
	AssertFloat("Acc", 216.558331574877, pp.Results.Acc, 0.0001, t)
	AssertFloat("Total", 639.5804375142908, pp.Results.Total, 0.0001, t)
}

func BenchmarkOsuSRSingle(b *testing.B) {
	osuFile, err := os.Open("Avenged Sevenfold - Save Me (Drummer) [Tragedy].osu")
	if err != nil {
		b.Fatal(err)
	}

	beatMap, err1 := beatmap.ParseFromReader(osuFile)
	if err1 != nil {
		b.Fatal(err)
	}

	for n := 0; n < b.N; n++ {
		osu.CalculateSingle(beatMap.HitObjects, beatMap.Difficulty)
	}
}

func BenchmarkOsuSRStep(b *testing.B) {
	osuFile, err := os.Open("Avenged Sevenfold - Save Me (Drummer) [Tragedy].osu")
	if err != nil {
		b.Fatal(err)
	}

	beatMap, err1 := beatmap.ParseFromReader(osuFile)
	if err1 != nil {
		b.Fatal(err)
	}

	for n := 0; n < b.N; n++ {
		osu.CalculateStep(beatMap.HitObjects, beatMap.Difficulty)
	}
}
