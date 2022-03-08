package archive_test

import (
	"archive/tar"
	"archive/zip"
	"errors"
	"fmt"
	"log"
	"testing"

	"github.com/kristinjeanna/archive"
)

type typeTest struct {
	filename      string
	expectedType  archive.Type
	expectedError error
}

var types = []typeTest{
	{"foo.tar", archive.Tar, nil},
	{"foo.tar.bz2", archive.TarBz2, nil},
	{"foo.tar.bzip2", archive.TarBz2, nil},
	{"foo.tbz", archive.TarBz2, nil},
	{"foo.tbz2", archive.TarBz2, nil},
	{"foo.tar.gz", archive.TarGz, nil},
	{"foo.tgz", archive.TarGz, nil},
	{"foo.tar.xz", archive.TarXz, nil},
	{"foo.txz", archive.TarXz, nil},
	{"foo.zip", archive.Zip, nil},
	{"foo.123", 0, archive.ErrUnknownType},
	{"foo.tar1", 0, archive.ErrUnknownType},
	{"foo.is.bar.abc.123-579.wxyz.tar.bz2", archive.TarBz2, nil},
	{"/usr/local/bin/foo.txz", archive.TarXz, nil},
	{`C:\Users\sam\Desktop\foo.zip`, archive.Zip, nil},
	{"foo.tAr.XZ", archive.TarXz, nil},
	{"foo.zip.tar.xz", archive.TarXz, nil},
	{"foo.tar.xz.zip", archive.Zip, nil},
}

func ExampleDetermineType() {
	filename := "sample.tar.gz"
	typ, err := archive.DetermineType(filename)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("File %s is type %s.", filename, typ)
	// Output: File sample.tar.gz is type TarGz.
}

func TestDetermineType(t *testing.T) {
	for _, c := range types {
		resultType, resultErr := archive.DetermineType(c.filename)

		if resultType != c.expectedType {
			t.Errorf("Expecting '%s', got '%s'\n", c.expectedType, resultType)
		}

		if resultErr != c.expectedError {
			t.Errorf("Expecting '%s', got '%s'\n", c.expectedError, resultErr)
		}
	}
}

func ExampleWalkTar() {
	callback := func(reader *tar.Reader, header *tar.Header) error {
		fmt.Printf("%s\n", header.Name)
		return nil
	}

	err := archive.WalkTar("testdata/sample.tar", callback)
	if err != nil {
		log.Fatal(err)
	}
	// Output:
	// sample/text/lorem.txt
	// sample/text/
	// sample/
}

func TestWalkTar(t *testing.T) {
	callback := func(reader *tar.Reader, header *tar.Header) error {
		fmt.Printf("%s\n", header.Name)
		return nil
	}

	err := archive.WalkTar("nonexistent.tar", callback)
	if err == nil {
		t.Error("Failed to receive non-nil error when walking a nonexistent tar file.")
	}

	err = archive.WalkTar("testdata/invalid.tar", callback)
	if err == nil {
		t.Error("Failed to receive non-nil error when walking an invalid tar file.")
	}
}

func ExampleWalkTarBzip2() {
	callback := func(reader *tar.Reader, header *tar.Header) error {
		fmt.Printf("%s\n", header.Name)
		return nil
	}

	err := archive.WalkTarBzip2("testdata/sample.tar.bz2", callback)
	if err != nil {
		log.Fatal(err)
	}
	// Output:
	// sample/
	// sample/text/
	// sample/text/lorem.txt
}

func TestWalkTarBzip2(t *testing.T) {
	callback := func(reader *tar.Reader, header *tar.Header) error {
		fmt.Printf("%s\n", header.Name)
		return nil
	}

	err := archive.WalkTarBzip2("nonexistent.tar.bz2", callback)
	if err == nil {
		t.Error("Failed to receive non-nil error when walking a nonexistent tar.bz2 file.")
	}
}

func ExampleWalkTarGz() {
	callback := func(reader *tar.Reader, header *tar.Header) error {
		fmt.Printf("%s\n", header.Name)
		return nil
	}

	err := archive.WalkTarGz("testdata/sample.tar.gz", callback)
	if err != nil {
		log.Fatal(err)
	}
	// Output:
	// sample/
	// sample/text/
	// sample/text/lorem.txt
}

func TestWalkTarGz(t *testing.T) {
	callback := func(reader *tar.Reader, header *tar.Header) error {
		fmt.Printf("%s\n", header.Name)
		return nil
	}

	err := archive.WalkTarGz("nonexistent.tar.gz", callback)
	if err == nil {
		t.Error("Failed to receive non-nil error when walking a nonexistent tar.gz file.")
	}

	err = archive.WalkTarGz("testdata/sample.tar.xz", callback)
	if err == nil {
		t.Error("Failed to receive non-nil error when walking a tar.xz file in WalkTarGz.")
	}
}

func ExampleWalkTarXz() {
	callback := func(reader *tar.Reader, header *tar.Header) error {
		fmt.Printf("%s\n", header.Name)
		return nil
	}

	err := archive.WalkTarXz("testdata/sample.tar.xz", callback)
	if err != nil {
		log.Fatal(err)
	}
	// Output:
	// sample/
	// sample/text/
	// sample/text/lorem.txt
}

func TestWalkTarXz(t *testing.T) {
	callback := func(reader *tar.Reader, header *tar.Header) error {
		fmt.Printf("%s\n", header.Name)
		return nil
	}

	err := archive.WalkTarXz("nonexistent.tar.xz", callback)
	if err == nil {
		t.Error("Failed to receive non-nil error when walking a nonexistent tar.xz file.")
	}

	err = archive.WalkTarXz("testdata/sample.tar.gz", callback)
	if err == nil {
		t.Error("Failed to receive non-nil error when walking a tar.gz file via the WalkTarXz function.")
	}

	callback = func(reader *tar.Reader, header *tar.Header) error {
		return errors.New("an error in callback processing")
	}

	err = archive.WalkTarXz("testdata/sample.tar.xz", callback)
	if err == nil {
		t.Error("Failed to return error from callback.")
	}
}

func ExampleWalkZip() {
	callback := func(file *zip.File) error {
		fmt.Printf("%s\n", file.Name)
		return nil
	}

	err := archive.WalkZip("testdata/sample.zip", callback)
	if err != nil {
		log.Fatal(err)
	}
	// Output:
	// sample/
	// sample/text/
	// sample/text/lorem.txt
}

func TestWalkzip(t *testing.T) {
	callback := func(file *zip.File) error {
		fmt.Printf("%s\n", file.Name)
		return nil
	}

	err := archive.WalkZip("nonexistent.zip", callback)
	if err == nil {
		t.Error("Failed to receive non-nil error when walking a nonexistent zip file.")
	}

	callback = func(file *zip.File) error {
		return errors.New("an error in callback processing")
	}

	err = archive.WalkZip("testdata/sample.zip", callback)
	if err == nil {
		t.Error("Failed to return error from callback.")
	}
}

type typeStringTest struct {
	archiveType archive.Type
	expected    string
}

var typeStrings = []typeStringTest{
	{archive.Tar, "Tar"},
	{archive.TarBz2, "TarBz2"},
	{archive.TarGz, "TarGz"},
	{archive.TarXz, "TarXz"},
	{archive.Zip, "Zip"},
}

func TestType_String(t *testing.T) {
	for _, typ := range typeStrings {
		result := typ.archiveType.String()

		if result != typ.expected {
			t.Errorf("Expecting '%s', got '%s'\n", typ.expected, result)
		}
	}
}
