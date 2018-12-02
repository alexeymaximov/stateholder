package stateholder

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

var testPath = filepath.Join(os.TempDir(), "test.mem")

var testBytes = []byte{'H', 'E', 'L', 'L', 'O'}
var testUint64 = uint64(1024)

var emptyBytes = []byte{0, 0, 0, 0, 0}
var emptyUint64 = uint64(0)

func clearStateholder() error {
	if _, err := os.Stat(testPath); err == nil || !os.IsNotExist(err) {
		if err := os.Remove(testPath); err != nil {
			return err
		}
	}
	return nil
}

func testStateholder() *Stateholder {
	stateholder := NewStateholder()
	stateholder.Define("bytes", uint16(len(testBytes)))
	stateholder.DefineUint64("uint64")
	return stateholder
}

func TestFirstOpen(t *testing.T) {
	if err := clearStateholder(); err != nil {
		t.Fatal(err)
	}
	stateholder := testStateholder()
	defer stateholder.Close()
	if _, err := stateholder.Attach(testPath, nil); err != nil {
		t.Fatal(err)
	}
	if value, err := stateholder.Get("bytes"); err != nil {
		t.Fatal(err)
	} else if bytes.Compare(value, emptyBytes) != 0 {
		t.Fatalf("bytes must be a %v, %v found", emptyBytes, value)
	}
	if value, err := stateholder.GetUint64("uint64"); err != nil {
		t.Fatal(err)
	} else if value != emptyUint64 {
		t.Fatalf("uint64 must be a %d, %d found", emptyUint64, value)
	}
}

func TestReadWrite(t *testing.T) {
	stateholder := testStateholder()
	defer stateholder.Close()
	if _, err := stateholder.Attach(testPath, nil); err != nil {
		t.Fatal(err)
	}
	if err := stateholder.Set("bytes", testBytes); err != nil {
		t.Fatal(err)
	}
	if err := stateholder.SetUint64("uint64", testUint64); err != nil {
		t.Fatal(err)
	}
	if value, err := stateholder.Get("bytes"); err != nil {
		t.Fatal(err)
	} else if bytes.Compare(value, testBytes) != 0 {
		t.Fatalf("bytes must be a %q, %v found", testBytes, value)
	}
	if value, err := stateholder.GetUint64("uint64"); err != nil {
		t.Fatal(err)
	} else if value != testUint64 {
		t.Fatalf("uint64 must be a %d, %d found", testUint64, value)
	}
}

func TestWrongOrder(t *testing.T) {
	stateholder := testStateholder()
	defer stateholder.Close()
	if _, err := stateholder.GetUint64("uint64"); err == nil {
		t.Fatal("expected ErrorDetached, no error found")
	} else if _, ok := err.(*ErrorDetached); !ok {
		t.Fatalf("expected ErrorDetached, [%v] error found", err)
	}
	if _, err := stateholder.Attach(testPath, nil); err != nil {
		t.Fatal(err)
	}
	if err := stateholder.DefineByte("byte"); err == nil {
		t.Fatal("expected ErrorAttached, no error found")
	} else if _, ok := err.(*ErrorAttached); !ok {
		t.Fatalf("expected ErrorAttached, [%v] error found", err)
	}
}

func TestWrongWrite(t *testing.T) {
	stateholder := testStateholder()
	defer stateholder.Close()
	if _, err := stateholder.Attach(testPath, nil); err != nil {
		t.Fatal(err)
	}
	if err := stateholder.SetByte("byte", 8); err == nil {
		t.Fatal("expected ErrorUndefined, no error found")
	} else if _, ok := err.(*ErrorUndefined); !ok {
		t.Fatalf("expected ErrorUndefined, [%v] error found", err)
	}
	if err := stateholder.Set("bytes", []byte{'H', 'E', 'L', 'L', 'O', '!'}); err == nil {
		t.Fatal("expected ErrorIncompatibleSize, no error found")
	} else if _, ok := err.(*ErrorIncompatibleSize); !ok {
		t.Fatalf("expected ErrorIncompatibleSize, [%v] error found", err)
	}
	if err := stateholder.SetUint32("uint64", 2048); err == nil {
		t.Fatal("expected ErrorIncompatibleKind, no error found")
	} else if _, ok := err.(*ErrorIncompatibleKind); !ok {
		t.Fatalf("expected ErrorIncompatibleKind, [%v] error found", err)
	}
}

func TestSubsequentRead(t *testing.T) {
	stateholder := testStateholder()
	defer stateholder.Close()
	if _, err := stateholder.Attach(testPath, nil); err != nil {
		t.Fatal(err)
	}
	if value, err := stateholder.Get("bytes"); err != nil {
		t.Fatal(err)
	} else if bytes.Compare(value, testBytes) != 0 {
		t.Fatalf("bytes must be a %q, %v found", testBytes, value)
	}
	if value, err := stateholder.GetUint64("uint64"); err != nil {
		t.Fatal(err)
	} else if value != testUint64 {
		t.Fatalf("uint64 must be a %d, %d found", testUint64, value)
	}
}

func TestTransactionCommit(t *testing.T) {
	if err := clearStateholder(); err != nil {
		t.Fatal(err)
	}
	stateholder := testStateholder()
	defer stateholder.Close()
	if _, err := stateholder.Attach(testPath, nil); err != nil {
		t.Fatal(err)
	}
	if err := stateholder.Begin(); err != nil {
		t.Fatal(err)
	}
	if err := stateholder.SetUint64("uint64", testUint64); err != nil {
		t.Fatal(err)
	}
	if value, err := stateholder.GetUint64("uint64"); err != nil {
		t.Fatal(err)
	} else if value != testUint64 {
		t.Fatalf("uint64 must be a %d, %d found", testUint64, value)
	}
	if err := stateholder.Commit(); err != nil {
		t.Fatal(err)
	}
	if value, err := stateholder.GetUint64("uint64"); err != nil {
		t.Fatal(err)
	} else if value != testUint64 {
		t.Fatalf("uint64 must be a %d, %d found", testUint64, value)
	}
}

func TestTransactionRollback(t *testing.T) {
	if err := clearStateholder(); err != nil {
		t.Fatal(err)
	}
	stateholder := testStateholder()
	defer stateholder.Close()
	if _, err := stateholder.Attach(testPath, nil); err != nil {
		t.Fatal(err)
	}
	if err := stateholder.Begin(); err != nil {
		t.Fatal(err)
	}
	if err := stateholder.SetUint64("uint64", testUint64); err != nil {
		t.Fatal(err)
	}
	if value, err := stateholder.GetUint64("uint64"); err != nil {
		t.Fatal(err)
	} else if value != testUint64 {
		t.Fatalf("uint64 must be a %d, %d found", testUint64, value)
	}
	if err := stateholder.Rollback(); err != nil {
		t.Fatal(err)
	}
	if value, err := stateholder.GetUint64("uint64"); err != nil {
		t.Fatal(err)
	} else if value != emptyUint64 {
		t.Fatalf("uint64 must be a %d, %d found", emptyUint64, value)
	}
}

func BenchmarkSync(b *testing.B) {
	stateholder := testStateholder()
	defer stateholder.Close()
	if _, err := stateholder.Attach(testPath, nil); err != nil {
		b.Fatal(err)
	}
	for i := 1; i <= 100; i++ {
		if err := stateholder.SetUint64("uint64", uint64(i)); err != nil {
			b.Fatal(err)
		}
		if err := stateholder.Sync(); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkNoSync(b *testing.B) {
	stateholder := testStateholder()
	defer stateholder.Close()
	if _, err := stateholder.Attach(testPath, nil); err != nil {
		b.Fatal(err)
	}
	for i := 1; i <= 100; i++ {
		if err := stateholder.SetUint64("uint64", uint64(i)); err != nil {
			b.Fatal(err)
		}
	}
}
