package reader

import (
	"archive/zip"
	"bytes"

	"github.com/spf13/afero/zipfs"
)

type ZipArchive struct {
	Archive
}

func (a *ZipArchive) Contents(b []byte) ([]byte, error) {
	rzip, err := zip.NewReader(bytes.NewReader(b), int64(len(b)))
	if err != nil {
		return nil, err
	}
	zipFS := zipfs.New(rzip)

	return a.Archive.FsContents(zipFS)
}
