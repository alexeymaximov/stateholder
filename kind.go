package stateholder

// Kind.
type Kind byte

// Available kinds.
const (
	KindBytes Kind = iota
	KindByte
	KindUint16
	KindUint32
	KindUint64
)

// Stringify kind.
func (kind Kind) String() string {
	switch kind {
	case KindBytes:
		return "byte array"
	case KindByte:
		return "byte"
	case KindUint16:
		return "uint16"
	case KindUint32:
		return "uint32"
	case KindUint64:
		return "uint64"
	default:
		return "invalid kind"
	}
}
