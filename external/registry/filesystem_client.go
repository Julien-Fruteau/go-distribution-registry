package registry

import (
	"bytes"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/julien-fruteau/go-distribution-registry/internal/env"
)

var fsBlobsPath = env.GetEnvOrDefault("REG_BLOBS_PATH", "/var/lib/registry/docker/registry/v2/blobs/sha256")

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

func WalkDirFn(path string, d fs.DirEntry, err error) error {
  if err != nil {
    return err
  }
  if d.IsDir() || len(path) <= len(fsBlobsPath)+4 {
    return nil
  }

  gzip, err := IsFileGzip(path)
  if err != nil {
    return err
  }
  if gzip {
    // TODO: something useful
    fmt.Println(path)
  }

	return nil
}

func WalkFs() error {
	err := filepath.WalkDir(fsBlobsPath, WalkDirFn)
	if err != nil {
		return err
	}
	return nil
}
