package test

import (
	"github.com/Wieku/gosu-pp/beatmap"
	"os"
	"testing"
)

func TestMetadata(t *testing.T) {
	osuFile, err := os.Open("inabakumori - Lost Umbrella (Ryuusei Aika) [________].osu")
	if err != nil {
		t.Fatal(err)
	}

	beatMap, err1 := beatmap.ParseFromReader(osuFile)
	if err1 != nil {
		t.Fatal(err)
	}

	Assert("Artist", "inabakumori", beatMap.Artist, t)
	Assert("ArtistUnicode", "稲葉曇", beatMap.ArtistUnicode, t)
	Assert("Title", "Lost Umbrella", beatMap.Title, t)
	Assert("TitleUnicode", "ロストアンブレラ", beatMap.TitleUnicode, t)
}

func TestMetadataUTF16(t *testing.T) {
	osuFile, err := os.Open("inabakumori - Lost Umbrella (Ryuusei Aika) [________] utf16.osu")
	if err != nil {
		t.Fatal(err)
	}

	beatMap, err1 := beatmap.ParseFromReader(osuFile)
	if err1 != nil {
		t.Fatal(err)
	}

	Assert("Artist", "inabakumori", beatMap.Artist, t)
	Assert("ArtistUnicode", "稲葉曇", beatMap.ArtistUnicode, t)
	Assert("Title", "Lost Umbrella", beatMap.Title, t)
	Assert("TitleUnicode", "ロストアンブレラ", beatMap.TitleUnicode, t)
}
