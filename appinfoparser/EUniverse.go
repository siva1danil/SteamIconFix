package appinfoparser

type EUniverse uint32

const (
	EUniverseInvalid  EUniverse = 0
	EUniversePublic   EUniverse = 1
	EUniverseBeta     EUniverse = 2
	EUniverseInternal EUniverse = 3
	EUniverseDev      EUniverse = 4
	EUniverseMax      EUniverse = 5
)
