package registry

import (
	"bytes"
	"io"
)

func ReadMagicBytes(r io.Reader) []byte {
	buf := make([]byte, 2)
	r.Read(buf)
	return buf
}

func IsGzipMagicBytes(b []byte) bool {
	return bytes.Equal(b, []byte{0x1F, 0x8B})
}

// TODO: 
// walk filesystem
// return blobs digest if gzip
