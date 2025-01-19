package appinfoparser

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"unicode/utf8"
)

const (
	Magic   uint32 = 0x07564427
	Magic28 uint32 = 0x07564428
	Magic29 uint32 = 0x07564429
)

type AppInfo struct {
	Universe    EUniverse
	Apps        []App
	StringTable []string
}

func AppInfoFromReader(r io.Reader) (*AppInfo, error) {
	info := &AppInfo{}
	tmp := uint32(0)
	reader := &CountingReader{r: r}

	if err := binary.Read(reader, binary.LittleEndian, &tmp); err != nil {
		return nil, err
	}
	if tmp != Magic && tmp != Magic28 && tmp != Magic29 {
		return nil, fmt.Errorf("unknown magic header: 0x%X", tmp)
	}

	if err := binary.Read(reader, binary.LittleEndian, &info.Universe); err != nil {
		return nil, err
	}

	if tmp == Magic29 {
		var stringTableOffset int64
		if err := binary.Read(reader, binary.LittleEndian, &stringTableOffset); err != nil {
			return nil, err
		}

		offset := reader.pos
		if err := reader.SeekRelative(int64(stringTableOffset - offset)); err != nil {
			return nil, err
		}

		var stringCount uint32
		if err := binary.Read(reader, binary.LittleEndian, &stringCount); err != nil {
			return nil, err
		}
		info.StringTable = make([]string, stringCount)

		for i := 0; i < int(stringCount); i++ {
			s, err := readNullTermUtf8String(reader)
			if err != nil {
				return nil, err
			}
			info.StringTable[i] = s
		}

		if err := reader.SeekRelative(offset + 0 - reader.pos); err != nil {
			return nil, err
		}
	}

	for {
		app := App{}

		if err := binary.Read(reader, binary.LittleEndian, &app.AppID); err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return nil, err
		}
		if app.AppID == 0 {
			break
		}

		size := uint32(0)
		if err := binary.Read(reader, binary.LittleEndian, &size); err != nil {
			return nil, err
		}
		endPos := reader.pos + int64(size)

		if err := binary.Read(reader, binary.LittleEndian, &app.InfoState); err != nil {
			return nil, err
		}
		if err := binary.Read(reader, binary.LittleEndian, &app.LastUpdated); err != nil {
			return nil, err
		}
		if err := binary.Read(reader, binary.LittleEndian, &app.Token); err != nil {
			return nil, err
		}
		app.Hash = make([]byte, 20)
		if _, err := io.ReadFull(reader, app.Hash); err != nil {
			return nil, err
		}
		if err := binary.Read(reader, binary.LittleEndian, &app.ChangeNumber); err != nil {
			return nil, err
		}
		if tmp == Magic28 || tmp == Magic29 {
			app.BinaryDataHash = make([]byte, 20)
			if _, err := io.ReadFull(reader, app.BinaryDataHash); err != nil {
				return nil, err
			}
		}

		dataSize := endPos - reader.pos
		if dataSize < 0 {
			return nil, errors.New("invalid data size")
		}
		data := make([]byte, dataSize)
		if _, err := io.ReadFull(reader, data); err != nil {
			return nil, err
		}
		kvdata, err := DataFromReader(info.StringTable, bytes.NewReader(data))
		if err != nil {
			return nil, err
		}
		app.Data = kvdata

		if reader.pos != endPos {
			return nil, errors.New("unexpected size mismatch after reading app data")
		}

		info.Apps = append(info.Apps, app)
	}

	return info, nil
}

func readNullTermUtf8String(r io.Reader) (string, error) {
	var buf bytes.Buffer
	tmp := make([]byte, 1)

	for {
		_, err := r.Read(tmp)
		if err != nil {
			return "", err
		}

		if tmp[0] == 0 {
			break
		}

		buf.WriteByte(tmp[0])
	}

	result := buf.Bytes()
	if !utf8.Valid(result) {
		return "", fmt.Errorf("invalid utf-8 data")
	}
	return string(result), nil
}
