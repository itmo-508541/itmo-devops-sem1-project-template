package main

import (
	"archive/tar"
	"archive/zip"
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/spf13/afero"
	"github.com/spf13/afero/tarfs"
	"github.com/spf13/afero/zipfs"
)

func main() {
}

func zf() {
	rzip, err := zip.OpenReader("build/test.zip")
	if err != nil {
		log.Fatal(err)
	}
	defer rzip.Close()
	// тут rzip.Reader - это zip.Reader

	// обязательно нужно как-то создать &zip.Reader
	zipFS := zipfs.New(&rzip.Reader)
	content, err := afero.ReadFile(zipFS, "/test.csv")
	if err != nil {
		panic(err)
	}

	fmt.Println(string(content))
}

func tf() {
	file, _ := os.Open("build/test.tar")
	reader := bufio.NewReader(file)
	rtar := tar.NewReader(reader)
	defer file.Close()

	// обязательно нужно как-то создать &tar.Reader
	tarFS := tarfs.New(rtar)
	content, _ := afero.ReadFile(tarFS, "/test.csv")

	fmt.Println(string(content))
}

func HttHandler(req *http.Request) {
	f, _, err := req.FormFile("fileTag")

	if err != nil {
		panic(err)
	}

	buf := new(bytes.Buffer)

	fileSize, err := io.Copy(buf, f)
	if err != nil {
		panic(err)
	}

	zip.NewReader(bytes.NewReader(buf.Bytes()), fileSize)

	// file, header, err := r.FormFile("zipfile")
	// // do something with error
	// zipReader, err := zip.NewReader(file, header.Size)
	// // do something with error
}
