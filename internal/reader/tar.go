package reader

import (
	"archive/tar"
	"bytes"
	"io"

	"github.com/spf13/afero/tarfs"
)

type TarArchive struct {
	Archive
}

func (a *TarArchive) Contents(b []byte) ([]byte, error) {
	return a.ReadContents(bytes.NewReader(b))
}

func (a *TarArchive) ReadContents(r io.Reader) ([]byte, error) {
	rtar := tar.NewReader(r)
	tarFS := tarfs.New(rtar)

	return a.FsContents(tarFS)
}
