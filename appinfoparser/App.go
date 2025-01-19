package appinfoparser

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
