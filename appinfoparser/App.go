package appinfoparser

// Partially based on code from SteamDB's AppInfo parser.
// Original code licensed under the MIT License:
// https://github.com/SteamDatabase/SteamAppInfo/blob/master/LICENSE
//
// Copyright (c) 2020 SteamDB

type App struct {
	AppID          uint32
	InfoState      uint32
	LastUpdated    uint32
	Token          uint64
	Hash           []byte
	BinaryDataHash []byte
	ChangeNumber   uint32
	Data           *Data
}
