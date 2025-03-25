package registry

import (
	"bytes"
	"compress/gzip"
	"math/rand/v2"
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
	if err != nil {
		t.Fatalf("failed to create file: %v", err)
	}
	defer f.Close()

	gzipWriter := gzip.NewWriter(f)
	defer gzipWriter.Close()

	_, err = gzipWriter.Write([]byte("Hello World"))
	if err != nil {
		t.Fatalf("failed to write to file: %v", err)
	}
}

func HelperCreateTxtFile(t testing.TB, path string) {
	t.Helper()
	f, err := os.Create(path)
	if err != nil {
		t.Fatalf("failed to create file: %v", err)
	}
	defer f.Close()

	_, err = f.Write([]byte("Hello World"))
	if err != nil {
		t.Fatalf("failed to write to file: %v", err)
	}
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

// ========== walk fs test
// needs at min 4 tests to validate the function

// used to create random directory name of desired len
func HelperRandomString(t testing.TB, n int) string {
	t.Helper()
	letters := []rune("abcdefghijklmnopqrstuvwxyz")

	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.IntN(len(letters))]
	}
	return string(b)
}

// create a random directory: root/<len2_intermediateDir>/<shaSumLen_finalDir>
// return the path of the directory
func HelperWalkFsCreateValidLenDir(t testing.TB, root string) string {
	t.Helper()
	subDir := filepath.Join(root, HelperRandomString(t, 2))
	err := os.Mkdir(subDir, 0777)
	if err != nil {
		t.Fatalf("failed to create dir: %v", err)
	}
	subDir2 := filepath.Join(subDir, HelperRandomString(t, sha256sumLen))
	err = os.Mkdir(subDir2, 0777)
	if err != nil {
		t.Fatalf("failed to create dir: %v", err)
	}
	return subDir2
}

func TestWalkFs_OnlyDirBadAndGodLenGetEmpty(t *testing.T) {
	walkDir := t.TempDir()
	// invalid length path
	err := os.Mkdir(filepath.Join(walkDir, "abc"), 0777)
	assert.NoError(t, err)
	HelperWalkFsCreateValidLenDir(t, walkDir)

	got, err := WalkFs(walkDir)

	assert.NoError(t, err)
	assert.Empty(t, got)
}

func TestWalkFs_WrongPathLenGzipFileGetEmpty(t *testing.T) {
	walkDir := t.TempDir()
	subDir := filepath.Join(walkDir, "tooShort")
	err := os.Mkdir(subDir, 0777)
	assert.NoError(t, err)
	// gzip file never reached by walkfs
	path := filepath.Join(subDir, "file.gz")
	HelperCreateGzipFile(t, path)

	got, err := WalkFs(walkDir)

	assert.NoError(t, err)
	assert.Empty(t, got)
}

func TestWalkFs_GoodPathLenTxtFileGetEmpty(t *testing.T) {
	walkDir := t.TempDir()
	subDir := HelperWalkFsCreateValidLenDir(t, walkDir)
	// text file not retained
	path := filepath.Join(subDir, "file.txt")
	HelperCreateTxtFile(t, path)

	got, err := WalkFs(walkDir)

	assert.NoError(t, err)
	assert.Empty(t, got)
}

func TestWalkFs_GoodPathLenGzipFileGetOne(t *testing.T) {
	walkDir := t.TempDir()
	subDir := HelperWalkFsCreateValidLenDir(t, walkDir)
	path := filepath.Join(subDir, "data")
	HelperCreateGzipFile(t, path)

	got, err := WalkFs(walkDir)

	assert.NoError(t, err)
	assert.Equal(t, 1, len(got))
	assert.Equal(t, path, got[0])
}

// ========== end walk fs test
