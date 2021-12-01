package beatmap

import (
	"bytes"
	"errors"
	"github.com/Wieku/gosu-pp/beatmap/objects"
	"github.com/Wieku/gosu-pp/files"
	"github.com/Wieku/gosu-pp/math/mutils"
	"io"
	"math"
	"sort"
	"strconv"
	"strings"
)

const bufferSize = 10 * 1024 * 1024

func parseGeneral(line []string, beatMap *BeatMap) bool {
	switch line[0] {
	case "Mode":
		beatMap.Mode, _ = strconv.Atoi(line[1])
	case "StackLeniency":
		beatMap.StackLeniency, _ = strconv.ParseFloat(line[1], 64)
		if math.IsNaN(beatMap.StackLeniency) {
			beatMap.StackLeniency = 0.0
		}
	case "AudioFilename":
		beatMap.Audio += line[1]
	case "PreviewTime":
		beatMap.PreviewTime, _ = strconv.ParseInt(line[1], 10, 64)
		//case "SampleSet":
		//	switch line[1] {
		//	case "Normal", "All":
		//		beatMap.Timings.BaseSet = 1
		//	case "Soft", "None":
		//		beatMap.Timings.BaseSet = 2
		//	case "Drum":
		//		beatMap.Timings.BaseSet = 3
		//	}
		//	beatMap.Timings.LastSet = beatMap.Timings.BaseSet
	}

	return false
}

func parseMetadata(line []string, beatMap *BeatMap) {
	switch line[0] {
	case "Title":
		beatMap.Title = line[1]
	case "TitleUnicode":
		beatMap.TitleUnicode = line[1]
	case "Artist":
		beatMap.Artist = line[1]
	case "ArtistUnicode":
		beatMap.ArtistUnicode = line[1]
	case "Creator":
		beatMap.Creator = line[1]
	case "FileVersion":
		beatMap.Version = line[1]
	case "Source":
		beatMap.Source = line[1]
	case "Tags":
		beatMap.Tags = line[1]
	case "BeatmapID":
		beatMap.MapID, _ = strconv.ParseInt(line[1], 10, 64)
	case "BeatmapSetID":
		beatMap.SetID, _ = strconv.ParseInt(line[1], 10, 64)
	}
}

func parseDifficulty(line []string, beatMap *BeatMap) {
	switch line[0] {
	case "SliderMultiplier":
		beatMap.SliderMultiplier, _ = strconv.ParseFloat(line[1], 64)
		beatMap.Timings.SliderMult = beatMap.SliderMultiplier
	case "ApproachRate":
		parsed, _ := strconv.ParseFloat(line[1], 64)
		beatMap.Difficulty.SetAR(mutils.ClampF64(parsed, 0, 10))
		beatMap.arSpecified = true
	case "CircleSize":
		parsed, _ := strconv.ParseFloat(line[1], 64)
		beatMap.Difficulty.SetCS(mutils.ClampF64(parsed, 0, 10))
	case "SliderTickRate":
		beatMap.Timings.TickRate, _ = strconv.ParseFloat(line[1], 64)
	case "HPDrainRate":
		parsed, _ := strconv.ParseFloat(line[1], 64)
		beatMap.Difficulty.SetHP(mutils.ClampF64(parsed, 0, 10))
	case "OverallDifficulty":
		parsed, _ := strconv.ParseFloat(line[1], 64)
		beatMap.Difficulty.SetOD(mutils.ClampF64(parsed, 0, 10))

		if !beatMap.arSpecified {
			beatMap.Difficulty.SetAR(beatMap.Difficulty.GetOD())
		}
	}
}

func parseEvents(line []string, beatMap *BeatMap) {
	switch line[0] {
	case "Background", "0":
		beatMap.Bg = strings.Replace(line[2], "\"", "", -1)
	case "Break", "2":
		beatMap.Pauses = append(beatMap.Pauses, NewPause(line))
	}
}

func parseHitObjects(line []string, beatMap *BeatMap) {
	obj := objects.CreateObject(line)

	if obj != nil {
		beatMap.HitObjects = append(beatMap.HitObjects, obj)
	}
}

func tokenize(line, delimiter string) []string {
	return tokenizeN(line, delimiter, -1)
}

func tokenizeN(line, delimiter string, n int) []string {
	if strings.HasPrefix(line, "//") || !strings.Contains(line, delimiter) {
		return nil
	}

	divided := strings.SplitN(line, delimiter, n)

	for i, a := range divided {
		divided[i] = strings.TrimSpace(a)
	}

	return divided
}

func getSection(line string) string {
	line = strings.TrimSpace(line)
	if strings.HasPrefix(line, "[") {
		return strings.TrimRight(strings.TrimLeft(line, "["), "]")
	}

	return ""
}

func ParseFromByte(data []byte) (*BeatMap, error) {
	return ParseFromReader(bytes.NewReader(data))
}

func ParseFromReader(reader io.Reader) (*BeatMap, error) {
	beatMap := NewBeatMap()

	scanner := files.NewScannerBuf(reader, bufferSize)

	var currentSection string

	counter := 0

	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(line, "osu file format v") {
			trim := strings.TrimPrefix(line, "osu file format v")
			beatMap.FileVersion, _ = strconv.Atoi(trim)
		}

		section := getSection(line)
		if section != "" {
			currentSection = section
			continue
		}

		switch currentSection {
		case "General":
			if arr := tokenizeN(line, ":", 2); len(arr) > 1 {
				parseGeneral(arr, beatMap)
			}
		case "Metadata":
			if arr := tokenizeN(line, ":", 2); len(arr) > 1 {
				parseMetadata(arr, beatMap)
			}
		case "Difficulty":
			if arr := tokenizeN(line, ":", 2); len(arr) > 1 {
				parseDifficulty(arr, beatMap)
			}
		case "Events":
			if arr := tokenize(line, ","); len(arr) > 1 {
				parseEvents(arr, beatMap)
			}
		case "TimingPoints":
			if arr := tokenize(line, ","); len(arr) > 1 {
				beatMap.ParsePoint(line)
				counter++
			}
		case "HitObjects":
			if arr := tokenize(line, ","); arr != nil {
				var time string

				objTypeI, _ := strconv.Atoi(arr[3])
				objType := objects.Type(objTypeI)
				if (objType & objects.CIRCLE) > 0 {
					beatMap.Circles++
					time = arr[2]
				} else if (objType & objects.SPINNER) > 0 {
					beatMap.Spinners++
					time = arr[5]
				} else if (objType & objects.SLIDER) > 0 {
					beatMap.Sliders++
					time = arr[2]
				} else if (objType & objects.LONGNOTE) > 0 {
					beatMap.Sliders++
					time = strings.Split(arr[5], ":")[0]
				}
				timeI, _ := strconv.Atoi(time)

				beatMap.Length = mutils.MaxI(beatMap.Length, timeI)

				parseHitObjects(arr, beatMap)
			}
		}
	}

	beatMap.FinalizePoints()

	if beatMap.Title+beatMap.Artist+beatMap.Creator == "" || counter == 0 {
		return nil, errors.New("corrupted file")
	}

	sort.SliceStable(beatMap.HitObjects, func(i, j int) bool {
		return beatMap.HitObjects[i].GetStartTime() < beatMap.HitObjects[j].GetStartTime()
	})

	num := 0
	comboNumber := 1
	comboSet := 0
	comboSetHax := 0
	forceNewCombo := false

	for _, iO := range beatMap.HitObjects {
		if iO.GetType() == objects.SPINNER {
			forceNewCombo = true
		} else if iO.IsNewCombo() || forceNewCombo {
			iO.SetNewCombo(true)
			comboNumber = 1
			comboSet++
			comboSetHax += int(iO.GetColorOffset()) + 1

			forceNewCombo = false
		}

		iO.SetID(num)
		iO.SetComboNumber(comboNumber)
		iO.SetComboSet(comboSet)
		iO.SetComboSetHax(comboSetHax)

		comboNumber++
		num++
	}

	for _, obj := range beatMap.HitObjects {
		obj.SetTiming(beatMap.Timings)
	}

	calculateStackLeniency(beatMap)

	return beatMap, nil
}
