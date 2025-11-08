package reader

import (
	"archive/tar"
	"bytes"

	"github.com/spf13/afero/tarfs"
)

type TarArchive struct {
	Archive
}

func (a *TarArchive) Contents(b []byte) ([]byte, error) {
	reader := bytes.NewReader(b)
	rtar := tar.NewReader(reader)
	tarFS := tarfs.New(rtar)

	return a.Archive.FsContents(tarFS)
}
