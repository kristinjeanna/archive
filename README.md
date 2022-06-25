# github.com/kristinjeanna/archive

[![GitHub license](https://img.shields.io/github/license/kristinjeanna/archive.svg?style=flat)](https://github.com/kristinjeanna/archive/blob/main/LICENSE) ![Last commit](https://img.shields.io/github/last-commit/kristinjeanna/archive?style=flat) ![Build and test package](https://github.com/kristinjeanna/archive/actions/workflows/build.yml/badge.svg?branch=main) [![Go Reference](https://pkg.go.dev/badge/github.com/kristinjeanna/archive.svg)](https://pkg.go.dev/github.com/kristinjeanna/archive)

Package `archive` is a convenience package for walking/enumerating the contents of zip files, tar files, and compressed tar files through callback functions. Supported archive types include: zip, tar, gzip-compressed tar, bzip2-compressed tar, and xz-compressed tar.

- [Install](#install)
- [Examples](#examples)
  - [List the contents of a .zip file](#list-the-contents-of-a-zip-file)
  - [Extract the contents of a .tar.xz file](#extract-the-contents-of-a-tarxz-file)
  - [Determine the type of archive file](#determine-the-type-of-archive-file)
- [Credits](#credits)

## Install

```shell
go get github.com/kristinjeanna/archive
```

## Examples

### List the contents of a .zip file

```go
func zipCallback(file *zip.File) error {
    if file.FileInfo().IsDir() {
        fmt.Printf("Dir : %s\n", file.Name)
    } else  {
        fmt.Printf("File: %s\n", file.Name)
    }

    return nil
}

func main() {
    err := archive.WalkZip("test.zip", zipCallback)
    if err != nil {
        log.Fatal(err)
    }
}

```

### Extract the contents of a .tar.xz file

```go
func tarCallback(reader *tar.Reader, header *tar.Header) error {
    if header.FileInfo().IsDir() {
        os.MkdirAll(header.Name, 0700)
        return nil
    }

    fo, err := os.Create(header.Name)
    if err != nil {
        return err
    }
    defer fo.Close()

    _, err := io.Copy(fo, reader)
    if err != nil {
        return err
    }

    return nil
}

func main() {
    err := archive.WalkTarXz("test.tar.xz", tarCallback)
    if err != nil {
        log.Fatal(err)
    }
}

```

### Determine the type of archive file

```go
func main() {
    archiveType, err := archive.DetermineType(archiveFilename)
    if err != nil {
        fmt.Fprintln(os.Stderr, "Unable to determine the file's archive type.")
        os.Exit(1)
    }

    switch archiveType {
    case archive.Tar:
        err = archive.WalkTar(archiveFilename, tarCallback)
    case archive.TarBz2:
        err = archive.WalkTarBzip2(archiveFilename, tarCallback)
    case archive.TarGz:
        err = archive.WalkTarGz(archiveFilename, tarCallback)
    case archive.TarXz:
        err = archive.WalkTarXz(archiveFilename, tarCallback)
    case archive.Zip:
        err = archive.WalkZip(archiveFilename, zipCallback)
    }

    if err != nil {
        log.Fatal(err)
    }
}

```

## Credits

- XZ compression support via [github.com/ulikunitz/xz](github.com/ulikunitz/xz)
