package reader

import (
	"fmt"
	"os"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

const testDataCsvHeader = "id,name,category,price,create_date"

func TestBadFilename(t *testing.T) {
	t.Parallel()

	osFs := afero.NewOsFs()
	testDataFs := afero.NewBasePathFs(osFs, "testdata/dir-not-exists")

	arch := &Archive{}
	filename, err := arch.Filename(testDataFs)

	assert.NotNil(t, err)
	assert.Equal(t, "", filename)
}

func TestNoFilename(t *testing.T) {
	t.Parallel()

	osFs := afero.NewOsFs()
	testDataFs := afero.NewBasePathFs(osFs, "testdata/empty")

	arch := &Archive{}
	filename, err := arch.Filename(testDataFs)

	assert.NotNil(t, err)
	assert.Equal(t, "", filename)
}

func TestFilename(t *testing.T) {
	t.Parallel()

	osFs := afero.NewOsFs()
	testDataFs := afero.NewBasePathFs(osFs, "testdata/dir")

	arch := &Archive{}
	filename, err := arch.Filename(testDataFs)

	assert.Nil(t, err)
	assert.Equal(t, "test.csv", filename)
}

func TestFsContents(t *testing.T) {
	t.Parallel()

	osFs := afero.NewOsFs()
	testDataFs := afero.NewBasePathFs(osFs, "testdata/dir")

	arch := &Archive{}
	contents, err := arch.FsContents(testDataFs)

	assert.Nil(t, err)
	assert.Equal(t, testDataCsvHeader, string(contents[:34]))
}

func TestContents(t *testing.T) {
	t.Parallel()

	tests := []struct {
		archive FileContents
		file    string
		title   string
		err     bool
	}{
		{archive: &ZipArchive{}, file: "testdata/zip/test.zip", title: "Zip", err: false},
		{archive: &ZipArchive{}, file: "testdata/tar/test.tar", title: "BadZip", err: true},
		{
			archive: &ZipArchive{},
			file:    "testdata/not-a-dir/not-a-file.txt",
			title:   "Zip NoFile",
			err:     true,
		},

		{archive: &TarArchive{}, file: "testdata/tar/test.tar", title: "Tar", err: false},
		{archive: &TarArchive{}, file: "testdata/zip/test.zip", title: "BadTar", err: true},
		{
			archive: &TarArchive{},
			file:    "testdata/not-a-dir/not-a-file.txt",
			title:   "Tar NoFile",
			err:     true,
		},

		{archive: &MultiArchive{}, file: "testdata/tar/test.tar", title: "Multi Tar", err: false},
		{archive: &MultiArchive{}, file: "testdata/zip/test.zip", title: "Multi Zip", err: false},
		{
			archive: &MultiArchive{},
			file:    "testdata/dir/test.csv",
			title:   "Multy Not Archive",
			err:     true,
		},
		{
			archive: &MultiArchive{},
			file:    "testdata/not-a-dir/not-a-file.txt",
			title:   "Multy NoFile",
			err:     true,
		},
	}
	for _, test := range tests {
		t.Run(fmt.Sprintf("%s(%s)", test.title, test.file), func(t *testing.T) {
			t.Parallel()

			data, _ := os.ReadFile(test.file)
			contents, err := test.archive.Contents(data)

			if test.err {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, testDataCsvHeader, string(contents[:34]))
			}
		})
	}
}
