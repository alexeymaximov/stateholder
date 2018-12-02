package stateholder

import "gitlab.studionx.ru/rnd/boots/0.1.1/lib/syspack"

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
