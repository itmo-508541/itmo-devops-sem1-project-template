package reader

import (
	"archive/zip"
	"bytes"
	"io"

	"github.com/spf13/afero/zipfs"
)

type ZipArchive struct {
	Archive
}

func (a *ZipArchive) Contents(b []byte) ([]byte, error) {
	return a.ReadContents(bytes.NewReader(b), int64(len(b)))
}

func (a *ZipArchive) ReadContents(r io.ReaderAt, size int64) ([]byte, error) {
	rzip, err := zip.NewReader(r, size)
	if err != nil {
		return nil, err
	}
	zipFS := zipfs.New(rzip)

	return a.FsContents(zipFS)
}
