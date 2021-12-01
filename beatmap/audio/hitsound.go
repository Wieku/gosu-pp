package audio

type HitSound int

const (
	Normal HitSound = 1 << iota
	Whistle
	Finish
	Clap
)
