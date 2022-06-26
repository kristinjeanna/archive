package archive

import (
	"archive/tar"
	"archive/zip"
	"errors"
	"fmt"
	"testing"

	"go.uber.org/goleak"
)

type typeTest struct {
	filename      string
	expectedType  Type
	expectedError error
}

var types = []typeTest{
	{"foo.tar", Tar, nil},
	{"foo.tar.bz2", TarBz2, nil},
	{"foo.tar.bzip2", TarBz2, nil},
	{"foo.tbz", TarBz2, nil},
	{"foo.tbz2", TarBz2, nil},
	{"foo.tar.gz", TarGz, nil},
	{"foo.tgz", TarGz, nil},
	{"foo.tar.xz", TarXz, nil},
	{"foo.txz", TarXz, nil},
	{"foo.zip", Zip, nil},
	{"foo.123", 0, errUnknownType},
	{"foo.tar1", 0, errUnknownType},
	{"foo.is.bar.abc.123-579.wxyz.tar.bz2", TarBz2, nil},
	{"/usr/local/bin/foo.txz", TarXz, nil},
	{`C:\Users\sam\Desktop\foo.zip`, Zip, nil},
	{"foo.tAr.XZ", TarXz, nil},
	{"foo.zip.tar.xz", TarXz, nil},
	{"foo.tar.xz.zip", Zip, nil},
}

func TestDetermineType(t *testing.T) {
	for _, c := range types {
		resultType, resultErr := DetermineType(c.filename)

		if resultType != c.expectedType {
			t.Errorf("Expecting '%s', got '%s'\n", c.expectedType, resultType)
		}

		if resultErr != c.expectedError {
			t.Errorf("Expecting '%s', got '%s'\n", c.expectedError, resultErr)
		}
	}
}

func TestWalkTar(t *testing.T) {
	callback := func(reader *tar.Reader, header *tar.Header) error {
		fmt.Printf("%s\n", header.Name)
		return nil
	}

	err := WalkTar("nonexistent.tar", callback)
	if err == nil {
		t.Error("Failed to receive non-nil error when walking a nonexistent tar file.")
	}

	err = WalkTar("testdata/invalid.tar", callback)
	if err == nil {
		t.Error("Failed to receive non-nil error when walking an invalid tar file.")
	}
}

func TestWalkTarBzip2(t *testing.T) {
	callback := func(reader *tar.Reader, header *tar.Header) error {
		fmt.Printf("%s\n", header.Name)
		return nil
	}

	err := WalkTarBzip2("nonexistent.tar.bz2", callback)
	if err == nil {
		t.Error("Failed to receive non-nil error when walking a nonexistent tar.bz2 file.")
	}
}

func TestWalkTarGz(t *testing.T) {
	callback := func(reader *tar.Reader, header *tar.Header) error {
		fmt.Printf("%s\n", header.Name)
		return nil
	}

	err := WalkTarGz("nonexistent.tar.gz", callback)
	if err == nil {
		t.Error("Failed to receive non-nil error when walking a nonexistent tar.gz file.")
	}

	err = WalkTarGz("testdata/sample.tar.xz", callback)
	if err == nil {
		t.Error("Failed to receive non-nil error when walking a tar.xz file in WalkTarGz.")
	}
}

func TestWalkTarXz(t *testing.T) {
	callback := func(reader *tar.Reader, header *tar.Header) error {
		fmt.Printf("%s\n", header.Name)
		return nil
	}

	err := WalkTarXz("nonexistent.tar.xz", callback)
	if err == nil {
		t.Error("Failed to receive non-nil error when walking a nonexistent tar.xz file.")
	}

	err = WalkTarXz("testdata/sample.tar.gz", callback)
	if err == nil {
		t.Error("Failed to receive non-nil error when walking a tar.gz file via the WalkTarXz function.")
	}

	callback = func(reader *tar.Reader, header *tar.Header) error {
		return errors.New("an error in callback processing")
	}

	err = WalkTarXz("testdata/sample.tar.xz", callback)
	if err == nil {
		t.Error("Failed to return error from callback.")
	}
}

func TestWalkzip(t *testing.T) {
	callback := func(file *zip.File) error {
		fmt.Printf("%s\n", file.Name)
		return nil
	}

	err := WalkZip("nonexistent.zip", callback)
	if err == nil {
		t.Error("Failed to receive non-nil error when walking a nonexistent zip file.")
	}

	callback = func(file *zip.File) error {
		return errors.New("an error in callback processing")
	}

	err = WalkZip("testdata/sample.zip", callback)
	if err == nil {
		t.Error("Failed to return error from callback.")
	}
}

type typeStringTest struct {
	archiveType Type
	expected    string
}

var typeStrings = []typeStringTest{
	{Tar, "Tar"},
	{TarBz2, "TarBz2"},
	{TarGz, "TarGz"},
	{TarXz, "TarXz"},
	{Zip, "Zip"},
}

func TestType_String(t *testing.T) {
	for _, typ := range typeStrings {
		result := typ.archiveType.String()

		if result != typ.expected {
			t.Errorf("Expecting '%s', got '%s'\n", typ.expected, result)
		}
	}
}

func TestMain(m *testing.M) {
	goleak.VerifyTestMain(m)
}
