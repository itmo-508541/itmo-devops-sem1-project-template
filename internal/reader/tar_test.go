package reader

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTarContents(t *testing.T) {
	data, _ := os.ReadFile("testdata/tar/test.tar")

	archive := &TarArchive{}
	contents, err := archive.Contents(data)

	assert.Nil(t, err)
	assert.Equal(t, "id,name,category,price,create_date", string(contents[:34]))
}
