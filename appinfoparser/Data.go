package appinfoparser

import "io"

type Data struct {
	Name        string
	MapValue    []Data
	StringValue *string
	NumberValue *uint32
}

func DataFromReader(dict []string, r io.Reader) (*Data, error) {
	return consume(dict, &CountingReader{r: r})
}

func consume(dict []string, r *CountingReader) (*Data, error) {
	typeByte := make([]byte, 1)
	_, err := r.Read(typeByte)
	if err != nil {
		return nil, err
	}

	if typeByte[0] == 0x00 {
		name, err := consumeNumber(r)
		if err != nil {
			return nil, err
		}
		value, err := consumeMap(dict, r)
		if err != nil {
			return nil, err
		}
		return &Data{Name: dict[name], MapValue: value}, nil
	} else if typeByte[0] == 0x01 {
		name, err := consumeNumber(r)
		if err != nil {
			return nil, err
		}
		value, err := consumeString(r)
		if err != nil {
			return nil, err
		}
		return &Data{Name: dict[name], StringValue: &value}, nil
	} else if typeByte[0] == 0x02 {
		name, err := consumeNumber(r)
		if err != nil {
			return nil, err
		}
		value, err := consumeNumber(r)
		if err != nil {
			return nil, err
		}
		return &Data{Name: dict[name], NumberValue: &value}, nil
	} else {
		return nil, nil
	}
}

func consumeString(r *CountingReader) (string, error) {
	s := ""
	b := make([]byte, 1)
	for {
		_, err := r.Read(b)
		if err != nil {
			return s, err
		}
		if b[0] == 0 {
			break
		}
		s += string(b[0])
	}
	return s, nil
}

func consumeNumber(r *CountingReader) (uint32, error) {
	b := make([]byte, 4)
	_, err := r.Read(b)
	if err != nil {
		return 0, err
	}
	return uint32(b[0]) | uint32(b[1])<<8 | uint32(b[2])<<16 | uint32(b[3])<<24, nil
}

func consumeMap(dict []string, r *CountingReader) ([]Data, error) {
	values := make([]Data, 0)
	for {
		value, err := consume(dict, r)
		if err != nil {
			return values, err
		}
		if value == nil {
			break
		}
		values = append(values, *value)
	}
	return values, nil
}
