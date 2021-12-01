package objects

import "strconv"

func CreateObject(data []string) IHitObject {
	objTypeI, _ := strconv.Atoi(data[3])
	objType := Type(objTypeI)

	if (objType & CIRCLE) > 0 {
		return NewCircle(data)
	} else if (objType & SPINNER) > 0 {
		return NewSpinner(data)
	} else if (objType & SLIDER) > 0 {
		return NewSlider(data)
	}

	return nil
}

type Type int

const (
	CIRCLE = Type(1 << iota)
	SLIDER
	NEWCOMBO
	SPINNER
	LONGNOTE = Type(128) //only for mania, used to have correct number of sliders in database just in case
)
