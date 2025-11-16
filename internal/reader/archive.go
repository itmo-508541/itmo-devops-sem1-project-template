package reader

import (
	"errors"
	"os"

	"github.com/spf13/afero"
)

type FileContents interface {
	Contents(b []byte) ([]byte, error)
}

type Archive struct{}

func (a *Archive) FsContents(fs afero.Fs) (result []byte, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New("Archive.FsContents: could not read file contents")
		}
	}()

	filename, err := a.Filename(fs)
	if err != nil {
		return nil, err
	}

	return afero.ReadFile(fs, filename)
}

func (a *Archive) Filename(fs afero.Fs) (fileName string, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.New("Archive.Filename: could not read from fs")
		}
	}()

	var files []os.FileInfo
	files, err = afero.ReadDir(fs, ".")
	if err != nil {
		return "", err
	}
	for _, file := range files {
		if file.Size() > 0 {
			return file.Name(), nil
		}
	}

	return "", errors.New("Archive.Filename: no files found")
}
