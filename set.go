package stateholder

import "encoding/binary"

// Set entry.
func (sh *Stateholder) set(key string, kind Kind, value []byte) error {
	if sh.index == nil {
		return &ErrorClosed{}
	}
	if sh.mapping == nil {
		return &ErrorDetached{}
	}
	index, ok := sh.index[key]
	if !ok {
		return &ErrorUndefined{Key: key}
	}
	entry := sh.entries[index]
	if kind != entry.kind {
		return &ErrorIncompatibleKind{Key: key, Kind: entry.kind, GivenKind: kind}
	}
	valueSize := EntrySize(len(value))
	if valueSize != entry.size {
		return &ErrorIncompatibleSize{Key: key, Size: entry.size, GivenSize: valueSize}
	}
	return sh.write(entry, value)
}

// Set byte array.
func (sh *Stateholder) Set(key string, value []byte) error {
	return sh.set(key, KindBytes, value)
}

// Set byte.
func (sh *Stateholder) SetByte(key string, value byte) error {
	return sh.set(key, KindByte, []byte{value})
}

// Set 16-bit unsigned integer value.
func (sh *Stateholder) SetUint16(key string, value uint16) error {
	buffer := make([]byte, 2)
	binary.LittleEndian.PutUint16(buffer, value)
	return sh.set(key, KindUint16, buffer)
}

// Set 32-bit unsigned integer value.
func (sh *Stateholder) SetUint32(key string, value uint32) error {
	buffer := make([]byte, 4)
	binary.LittleEndian.PutUint32(buffer, value)
	return sh.set(key, KindUint32, buffer)
}

// Set 64-bit unsigned integer value.
func (sh *Stateholder) SetUint64(key string, value uint64) error {
	buffer := make([]byte, 8)
	binary.LittleEndian.PutUint64(buffer, value)
	return sh.set(key, KindUint64, buffer)
}
