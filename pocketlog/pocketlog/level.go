package pocketlog

type Level uint8

const (
	LevelDebug Level = iota
	LevelInfo
	LevelWarning
	LevelError
)
