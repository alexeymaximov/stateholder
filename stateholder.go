package stateholder

// TODO: Use direct mapping slice for integer entries reading.

import (
	"bytes"
	"encoding/binary"
	"os"
	"runtime"

	"github.com/alexeymaximov/syspack"
	"github.com/alexeymaximov/syspack/mmap"
)

type Stateholder struct {
	// Stateholder.

	// Index.
	index map[string]int

	// Entries.
	entries []*entry

	// Size.
	size syspack.Size

	// Mapping.
	mapping *mmap.Mapping

	// Transaction mode.
	transaction bool
}

// Make new stateholder.
func NewStateholder() *Stateholder {
	sh := &Stateholder{index: make(map[string]int)}
	runtime.SetFinalizer(sh, (*Stateholder).Close)
	return sh
}

// Read entry.
func (sh *Stateholder) read(entry *entry) ([]byte, error) {
	value := make([]byte, entry.size)
	if sh.transaction && entry.buffer != nil {
		copy(value, entry.buffer)
	} else if n, err := sh.mapping.ReadAt(value, entry.offset); err != nil {
		return nil, err
	} else if n != int(entry.size) {
		return nil, &ErrorCorruptedRead{Real: n, Expected: int(entry.size)}
	}
	return value, nil
}

// Write entry.
func (sh *Stateholder) write(entry *entry, value []byte) error {
	if sh.transaction {
		if entry.buffer == nil {
			entry.buffer = make([]byte, entry.size)
		}
		copy(entry.buffer, value)
	} else if n, err := sh.mapping.WriteAt(value, entry.offset); err != nil {
		return err
	} else if n != int(entry.size) {
		return &ErrorCorruptedWrite{Real: n, Expected: int(entry.size)}
	}
	return nil
}

// Copy entry.
func (sh *Stateholder) Copy(key, sourceKey string) error {
	if sh.index == nil {
		return &ErrorClosed{}
	}
	if sh.mapping == nil {
		return &ErrorDetached{}
	}
	var index int
	var ok bool
	index, ok = sh.index[key]
	if !ok {
		return &ErrorUndefined{Key: key}
	}
	entry := sh.entries[index]
	index, ok = sh.index[sourceKey]
	if !ok {
		return &ErrorUndefined{Key: sourceKey}
	}
	sourceEntry := sh.entries[index]
	if entry.kind != sourceEntry.kind {
		return &ErrorIncompatibleKind{Key: key, Kind: entry.kind, GivenKind: sourceEntry.kind}
	}
	if entry.size != sourceEntry.size {
		return &ErrorIncompatibleSize{Key: key, Size: entry.size, GivenSize: sourceEntry.size}
	}
	value, err := sh.read(sourceEntry)
	if err != nil {
		return err
	}
	return sh.write(entry, value)
}

// Prepare file.
func (sh *Stateholder) prepareFile(file *os.File, sign []byte) error {
	signLen := len(sign)
	if err := file.Truncate(int64(syspack.Size(signLen) + sh.size)); err != nil {
		return err
	}
	if signLen > 0 {
		buffer := make([]byte, syspack.Size(signLen))
		copy(buffer[:signLen], sign)
		if n, err := file.WriteAt(buffer, 0); err != nil {
			return err
		} else if n != signLen {
			return &ErrorBadFile{Path: file.Name()}
		}
	}
	if err := file.Sync(); err != nil {
		return err
	}
	return nil
}

// Attach file and return true is new file was created.
func (sh *Stateholder) Attach(filePath string, sign []byte) (bool, error) {
	if sh.index == nil {
		return false, &ErrorClosed{}
	}
	if sh.mapping != nil {
		return false, &ErrorAttached{}
	}
	if sign == nil {
		sign = []byte{'M', 'E', 'M', 1, 0, 0}
		entrySign := make([]byte, 3)
		for _, entry := range sh.entries {
			entrySign[0] = byte(entry.kind)
			binary.LittleEndian.PutUint16(entrySign[1:], entry.size)
			sign = append(sign, entrySign...)
		}
	}
	init := false
	if _, err := os.Stat(filePath); err != nil && os.IsNotExist(err) {
		init = true
	}
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_RDWR, 0600)
	if err != nil {
		return false, err
	}
	defer file.Close()
	if init {
		if err := sh.prepareFile(file, sign); err != nil {
			return false, err
		}
	}
	signLen := len(sign)
	buffer := make([]byte, signLen)
	if n, err := file.ReadAt(buffer, 0); err != nil {
		return false, err
	} else if n != signLen {
		return false, &ErrorBadFile{Path: filePath}
	}
	if bytes.Compare(buffer, sign) != 0 {
		return false, &ErrorBadFile{Path: filePath}
	}
	mapping, err := mmap.NewMapping(file.Fd(), syspack.Offset(signLen), sh.size, &mmap.Options{
		Mode: mmap.ModeReadWrite,
	})
	if err != nil {
		return false, err
	}
	sh.mapping = mapping
	return init, nil
}

// Sync data.
func (sh *Stateholder) Sync() error {
	if sh.index == nil {
		return &ErrorClosed{}
	}
	if sh.mapping == nil {
		return &ErrorDetached{}
	}
	return sh.mapping.Sync()
}

// Close stateholder.
func (sh *Stateholder) Close() error {
	if sh.index == nil {
		return &ErrorClosed{}
	}
	if sh.mapping != nil {
		if err := sh.mapping.Close(); err != nil {
			return err
		}
		sh.mapping = nil
	}
	sh.index = nil
	sh.entries = nil
	runtime.SetFinalizer(sh, nil)
	return nil
}
