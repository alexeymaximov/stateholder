package stateholder

import "gitlab.studionx.ru/rnd/boots/0.1.1/lib/syspack"

// Define entry.
func (sh *Stateholder) define(key string, kind Kind, size EntrySize) error {
	if sh.index == nil {
		return &ErrorClosed{}
	}
	if sh.mapping != nil {
		return &ErrorAttached{}
	}
	if _, ok := sh.index[key]; ok {
		return &ErrorAmbiguous{Key: key}
	}
	if size <= 0 {
		return &ErrorInvalidSize{Key: key, Size: size}
	}
	sh.index[key] = len(sh.entries)
	sh.entries = append(sh.entries, &entry{kind: kind, offset: syspack.Offset(sh.size), size: size})
	sh.size += syspack.Size(size)
	return nil
}

// Define byte array.
func (sh *Stateholder) Define(key string, size EntrySize) error {
	return sh.define(key, KindBytes, size)
}

// Define byte.
func (sh *Stateholder) DefineByte(key string) error {
	return sh.define(key, KindByte, 1)
}

// Define 16-bit unsigned integer value.
func (sh *Stateholder) DefineUint16(key string) error {
	return sh.define(key, KindUint16, 2)
}

// Define 32-bit unsigned integer value.
func (sh *Stateholder) DefineUint32(key string) error {
	return sh.define(key, KindUint32, 4)
}

// Define 64-bit unsigned integer value.
func (sh *Stateholder) DefineUint64(key string) error {
	return sh.define(key, KindUint64, 8)
}
