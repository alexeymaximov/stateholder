package stateholder

import "fmt"

// Error occurred when key is ambiguous.
type ErrorAmbiguous struct{ Key string }

// Get error message.
func (err *ErrorAmbiguous) Error() string {
	return fmt.Sprintf("stateholder: ambiguous %q", err.Key)
}

// Error occurred when file attached.
type ErrorAttached struct{}

// Get error message.
func (err *ErrorAttached) Error() string {
	return "stateholder: file attached"
}

// Error occurred when file is bad.
type ErrorBadFile struct{ Path string }

// Get error message.
func (err *ErrorBadFile) Error() string {
	return fmt.Sprintf("stateholder: bad file %s", err.Path)
}

// Error occurred when stateholder closed.
type ErrorClosed struct{}

// Get error message.
func (err *ErrorClosed) Error() string {
	return "stateholder: closed"
}

// Error occurred on read corruption.
type ErrorCorruptedRead struct{ Real, Expected int }

// Get error message.
func (err *ErrorCorruptedRead) Error() string {
	return fmt.Sprintf("stateholder: read %d bytes instead of %d", err.Real, err.Expected)
}

// Error occurred on write corruption.
type ErrorCorruptedWrite struct{ Real, Expected int }

// Get error message.
func (err *ErrorCorruptedWrite) Error() string {
	return fmt.Sprintf("stateholder: write %d bytes instead of %d", err.Real, err.Expected)
}

// Error occurred when file detached.
type ErrorDetached struct{}

// Get error message.
func (err *ErrorDetached) Error() string {
	return "stateholder: file detached"
}

// Error occurred when entry kind is incompatible with given one.
type ErrorIncompatibleKind struct {
	Key       string
	Kind      Kind
	GivenKind Kind
}

// Get error message.
func (err *ErrorIncompatibleKind) Error() string {
	return fmt.Sprintf("stateholder: %q is %s, not %s", err.Key, err.Kind, err.GivenKind)
}

// Error occurred when entry size is incompatible with given one.
type ErrorIncompatibleSize struct {
	Key       string
	Size      uint16
	GivenSize uint16
}

// Get error message.
func (err *ErrorIncompatibleSize) Error() string {
	return fmt.Sprintf("stateholder: %q size is %d bytes, not %d", err.Key, err.Size, err.GivenSize)
}

// Error occurred when entry size is invalid.
type ErrorInvalidSize struct {
	Key  string
	Size uint16
}

// Get error message.
func (err *ErrorInvalidSize) Error() string {
	return fmt.Sprintf("stateholder: size %d of %q is invalid", err.Size, err.Key)
}

// Error occurred when transaction not started.
type ErrorTransactionNotStarted struct{}

// Get error message.
func (err *ErrorTransactionNotStarted) Error() string {
	return "stateholder: transaction not started"
}

// Error occurred when transaction already started.
type ErrorTransactionAlreadyStarted struct{}

// Get error message.
func (err *ErrorTransactionAlreadyStarted) Error() string {
	return "stateholder: transaction already started"
}

// Error occurred when key is undefined.
type ErrorUndefined struct{ Key string }

// Get error message.
func (err *ErrorUndefined) Error() string {
	return fmt.Sprintf("stateholder: undefined %q", err.Key)
}
