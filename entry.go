package stateholder

import "github.com/alexeymaximov/syspack"

// Entry size.
type EntrySize = uint16

type entry struct {
	// Entry.

	// Kind.
	kind Kind

	// Offset.
	offset syspack.Offset

	// Size.
	size EntrySize

	// Buffer.
	buffer []byte
}
