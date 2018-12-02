package stateholder

import "encoding/binary"

// Get entry.
func (sh *Stateholder) get(key string, kind Kind) (*entry, []byte, error) {
	if sh.index == nil {
		return nil, nil, &ErrorClosed{}
	}
	if sh.mapping == nil {
		return nil, nil, &ErrorDetached{}
	}
	index, ok := sh.index[key]
	if !ok {
		return nil, nil, &ErrorUndefined{Key: key}
	}
	entry := sh.entries[index]
	if kind != entry.kind {
		return nil, nil, &ErrorIncompatibleKind{Key: key, Kind: entry.kind, GivenKind: kind}
	}
	value, err := sh.read(entry)
	if err != nil {
		return nil, nil, err
	}
	return entry, value, nil
}

// Get byte array.
func (sh *Stateholder) Get(key string) ([]byte, error) {
	_, value, err := sh.get(key, KindBytes)
	if err != nil {
		return nil, err
	}
	return value, nil
}

// Get byte.
func (sh *Stateholder) GetByte(key string) (byte, error) {
	_, value, err := sh.get(key, KindByte)
	if err != nil {
		return 0, err
	}
	return value[0], nil
}

// Get 16-bit unsigned integer value.
func (sh *Stateholder) GetUint16(key string) (uint16, error) {
	_, value, err := sh.get(key, KindUint16)
	if err != nil {
		return 0, err
	}
	return binary.LittleEndian.Uint16(value), nil
}

// Get 32-bit unsigned integer value.
func (sh *Stateholder) GetUint32(key string) (uint32, error) {
	_, value, err := sh.get(key, KindUint32)
	if err != nil {
		return 0, err
	}
	return binary.LittleEndian.Uint32(value), nil
}

// Get 64-bit unsigned integer value.
func (sh *Stateholder) GetUint64(key string) (uint64, error) {
	_, value, err := sh.get(key, KindUint64)
	if err != nil {
		return 0, err
	}
	return binary.LittleEndian.Uint64(value), nil
}
