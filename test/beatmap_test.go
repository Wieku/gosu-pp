package test

import (
	"github.com/wieku/gosu-pp/beatmap"
	"os"
	"testing"
)

func TestMetadata(t *testing.T) {
	osuFile, err := os.Open("inabakumori - Lost Umbrella (Ryuusei Aika) [________].osu")
	if err != nil {
		t.Fatal(err)
	}

	beatMap, err1 := beatmap.ParseBeatMapFile(osuFile)
	if err1 != nil {
		t.Fatal(err)
	}

	Assert("inabakumori", beatMap.Artist, t)
	Assert("稲葉曇", beatMap.ArtistUnicode, t)
	Assert("Lost Umbrella", beatMap.Title, t)
	Assert("ロストアンブレラ", beatMap.TitleUnicode, t)
}

func TestMetadataUTF16(t *testing.T) {
	osuFile, err := os.Open("inabakumori - Lost Umbrella (Ryuusei Aika) [________] utf16.osu")
	if err != nil {
		t.Fatal(err)
	}

	beatMap, err1 := beatmap.ParseBeatMapFile(osuFile)
	if err1 != nil {
		t.Fatal(err)
	}

	Assert("inabakumori", beatMap.Artist, t)
	Assert("稲葉曇", beatMap.ArtistUnicode, t)
	Assert("Lost Umbrella", beatMap.Title, t)
	Assert("ロストアンブレラ", beatMap.TitleUnicode, t)
}
