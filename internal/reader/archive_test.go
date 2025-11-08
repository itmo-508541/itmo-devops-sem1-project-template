package reader

import (
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestBadFilename(t *testing.T) {
	osFs := afero.NewOsFs()
	testDataFs := afero.NewBasePathFs(osFs, "testdata/dir-not-exists")

	arch := &Archive{}
	filename, err := arch.Filename(testDataFs)

	assert.NotNil(t, err)
	assert.Equal(t, "", filename)
}

func TestNoFilename(t *testing.T) {
	osFs := afero.NewOsFs()
	testDataFs := afero.NewBasePathFs(osFs, "testdata/empty")

	arch := &Archive{}
	filename, err := arch.Filename(testDataFs)

	assert.NotNil(t, err)
	assert.Equal(t, "", filename)
}

func TestFilename(t *testing.T) {
	osFs := afero.NewOsFs()
	testDataFs := afero.NewBasePathFs(osFs, "testdata/dir")

	arch := &Archive{}
	filename, err := arch.Filename(testDataFs)

	assert.Nil(t, err)
	assert.Equal(t, "test.csv", filename)
}

func TestFsContents(t *testing.T) {
	osFs := afero.NewOsFs()
	testDataFs := afero.NewBasePathFs(osFs, "testdata/dir")

	arch := &Archive{}
	contents, err := arch.FsContents(testDataFs)

	assert.Nil(t, err)
	assert.Equal(t, "id,name,category,price,create_date", string(contents[:34]))
}
