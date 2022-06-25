package archive

import (
	"archive/tar"
	"archive/zip"
	"fmt"
	"log"
)

func ExampleDetermineType() {
	filename := "sample.tar.gz"
	typ, err := DetermineType(filename)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("File %s is type %s.", filename, typ)
	// Output: File sample.tar.gz is type TarGz.
}

func ExampleWalkTar() {
	callback := func(reader *tar.Reader, header *tar.Header) error {
		fmt.Printf("%s\n", header.Name)
		return nil
	}

	err := WalkTar("testdata/sample.tar", callback)
	if err != nil {
		log.Fatal(err)
	}
	// Output:
	// sample/text/lorem.txt
	// sample/text/
	// sample/
}

func ExampleWalkTarBzip2() {
	callback := func(reader *tar.Reader, header *tar.Header) error {
		fmt.Printf("%s\n", header.Name)
		return nil
	}

	err := WalkTarBzip2("testdata/sample.tar.bz2", callback)
	if err != nil {
		log.Fatal(err)
	}
	// Output:
	// sample/
	// sample/text/
	// sample/text/lorem.txt
}

func ExampleWalkTarGz() {
	callback := func(reader *tar.Reader, header *tar.Header) error {
		fmt.Printf("%s\n", header.Name)
		return nil
	}

	err := WalkTarGz("testdata/sample.tar.gz", callback)
	if err != nil {
		log.Fatal(err)
	}
	// Output:
	// sample/
	// sample/text/
	// sample/text/lorem.txt
}

func ExampleWalkTarXz() {
	callback := func(reader *tar.Reader, header *tar.Header) error {
		fmt.Printf("%s\n", header.Name)
		return nil
	}

	err := WalkTarXz("testdata/sample.tar.xz", callback)
	if err != nil {
		log.Fatal(err)
	}
	// Output:
	// sample/
	// sample/text/
	// sample/text/lorem.txt
}

func ExampleWalkZip() {
	callback := func(file *zip.File) error {
		fmt.Printf("%s\n", file.Name)
		return nil
	}

	err := WalkZip("testdata/sample.zip", callback)
	if err != nil {
		log.Fatal(err)
	}
	// Output:
	// sample/
	// sample/text/
	// sample/text/lorem.txt
}
