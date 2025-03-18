package registry

import (
	"bytes"
	"compress/gzip"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

// https://en.wikipedia.org/wiki/List_of_file_signatures
func TestHexDecNotations(t *testing.T) {
	b := &bytes.Buffer{}
	// decimal notation
	b.Write([]byte{31, 139})
	got := b.Bytes()

	// hexadecimal notation
	want := []byte{0x1F, 0x8B}
	assert.Equal(t, want, got)
}

func TestMagicBytes(t *testing.T) {
	var b bytes.Buffer
	b.Write([]byte{31, 139})

	got := ReadMagicBytes(&b)

	want := []byte{0x1F, 0x8B}
	assert.Equal(t, want, got)
}

func TestIsGzipMagicBytes(t *testing.T) {
	want := []byte{0x1F, 0x8B}
	assert.True(t, IsGzipMagicBytes(want))
}

func HelperCreateGzipFile(t testing.TB, path string) {
	t.Helper()
	f, err := os.Create(path)
	assert.NoError(t, err)
	defer f.Close()

	gzipWriter := gzip.NewWriter(f)
	defer gzipWriter.Close()

	_, err = gzipWriter.Write([]byte("Hello World"))
	assert.NoError(t, err)
}

func HelperCreateTxtFile(t testing.TB, path string) {
	t.Helper()
	f, err := os.Create(path)
	assert.NoError(t, err)
	defer f.Close()

	_, err = f.Write([]byte("Hello World"))
	assert.NoError(t, err)
}

func TestIsGzipFileNotFound(t *testing.T) {
	_, err := IsFileGzip("not-a-path")
	assert.ErrorIs(t, err, os.ErrNotExist)
}

func TestIsGzipFileNo(t *testing.T) {
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "test")
	HelperCreateTxtFile(t, path)

	got, err := IsFileGzip(path)
	assert.NoError(t, err)
	assert.False(t, got)
}

func TestIsGzipFileYes(t *testing.T) {
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "test.gz")
	HelperCreateGzipFile(t, path)

	got, err := IsFileGzip(path)
	assert.NoError(t, err)
	assert.True(t, got)
}

