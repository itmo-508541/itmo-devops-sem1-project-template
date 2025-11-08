package reader

import (
	"fmt"

	"github.com/spf13/afero"
)

type FileContents interface {
	Contents(b []byte) ([]byte, error)
}

type Archive struct {
}

func (a *Archive) FsContents(fs afero.Fs) ([]byte, error) {
	filename, err := a.Filename(fs)
	if err != nil {
		return nil, err
	}

	return afero.ReadFile(fs, filename)
}

func (a *Archive) Filename(fs afero.Fs) (string, error) {
	files, err := afero.ReadDir(fs, ".")
	fmt.Println("FILES:", files)
	if err != nil {
		return "", err
	}
	for _, file := range files {
		if file.Size() > 0 {
			return file.Name(), nil
		}
	}

	return "", fmt.Errorf("Archive.Filename: no files found")
}
