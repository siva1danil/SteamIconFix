package appinfoparser

// Partially based on code from SteamDB's AppInfo parser.
// Original code licensed under the MIT License:
// https://github.com/SteamDatabase/SteamDatabase/blob/master/LICENSE
//
// Copyright (c) 2020 SteamDB

type EUniverse uint32

const (
	EUniverseInvalid  EUniverse = 0
	EUniversePublic   EUniverse = 1
	EUniverseBeta     EUniverse = 2
	EUniverseInternal EUniverse = 3
	EUniverseDev      EUniverse = 4
	EUniverseMax      EUniverse = 5
)
