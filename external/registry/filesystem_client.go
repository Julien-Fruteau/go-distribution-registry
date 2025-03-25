package registry

import (
	"bytes"
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/julien-fruteau/go-distribution-registry/internal/env"
)

var (
	fsBlobsPath  = env.GetEnvOrDefault("REG_BLOBS_PATH", "/var/lib/registry/docker/registry/v2/blobs/sha256")
	sha256sumLen = 64
	sha386sumLen = 96
	SHA512sumLen = 128
	shaSumLen    = env.GetEnvOrDefaultInt("REG_SHASUM_LEN", sha256sumLen)
)

func IsFileGzip(name string) (bool, error) {
	f, err := os.Open(name)
	if err != nil {
		return false, err
	}
	defer f.Close()
	mb := ReadMagicBytes(f)
	return IsGzipMagicBytes(mb), nil
}

func ReadMagicBytes(r io.Reader) []byte {
	buf := make([]byte, 2)
	r.Read(buf)
	return buf
}

func IsGzipMagicBytes(b []byte) bool {
	return bytes.Equal(b, []byte{0x1F, 0x8B})
}

// lookup for gzip files in root/<len2_intermediateDir>/<shaSumLen_finalDir>/ directory
// check the directory len prior checking if the file is gzip
// NOTE: `/<len2_intermediateDir>/` counts for the len 4
func WalkDirFnGzipBlobs(root, path string, d fs.DirEntry, err error, gzipBlobs *[]string) error {
	if err != nil {
		return err
	}
	filepath.Dir(path)
	if d.IsDir() || len(filepath.Dir(path)) != len(root)+4+shaSumLen {
		return nil
	}

	gzip, err := IsFileGzip(path)
	if err != nil {
		return err
	}

	if gzip {
		*gzipBlobs = append(*gzipBlobs, path)
	}

	return nil
}

// walks root directory to find gzip blobs file data
// returns the list of file path found
func WalkFs(root string) ([]string, error) {
	var gzipBlobs []string

	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		return WalkDirFnGzipBlobs(root, path, d, err, &gzipBlobs)
	})
	if err != nil {
		return nil, err
	}

	return gzipBlobs, nil
}

// TODO: process walkfs to get the digest name only
