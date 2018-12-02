package stateholder

import "encoding/binary"

// Increment byte.
func (sh *Stateholder) IncByte(key string, delta byte) (byte, byte, error) {
	entry, buffer, err := sh.get(key, KindByte)
	if err != nil {
		return 0, 0, err
	}
	old := buffer[0]
	buffer[0] += delta
	if err := sh.write(entry, buffer); err != nil {
		return 0, 0, err
	}
	return old, buffer[0], nil
}

// Increment 16-bit unsigned integer value.
func (sh *Stateholder) IncUint16(key string, delta uint16) (uint16, uint16, error) {
	entry, buffer, err := sh.get(key, KindUint16)
	if err != nil {
		return 0, 0, err
	}
	value := binary.LittleEndian.Uint16(buffer)
	old := value
	value += delta
	binary.LittleEndian.PutUint16(buffer, value)
	if err := sh.write(entry, buffer); err != nil {
		return 0, 0, err
	}
	return old, value, nil
}

// Increment to 32-bit unsigned integer value.
func (sh *Stateholder) IncUint32(key string, delta uint32) (uint32, uint32, error) {
	entry, buffer, err := sh.get(key, KindUint32)
	if err != nil {
		return 0, 0, err
	}
	value := binary.LittleEndian.Uint32(buffer)
	old := value
	value += delta
	binary.LittleEndian.PutUint32(buffer, value)
	if err := sh.write(entry, buffer); err != nil {
		return 0, 0, err
	}
	return old, value, nil
}

// Increment 64-bit unsigned integer value.
func (sh *Stateholder) IncUint64(key string, delta uint64) (uint64, uint64, error) {
	entry, buffer, err := sh.get(key, KindUint64)
	if err != nil {
		return 0, 0, err
	}
	value := binary.LittleEndian.Uint64(buffer)
	old := value
	value += delta
	binary.LittleEndian.PutUint64(buffer, value)
	if err := sh.write(entry, buffer); err != nil {
		return 0, 0, err
	}
	return old, value, nil
}
