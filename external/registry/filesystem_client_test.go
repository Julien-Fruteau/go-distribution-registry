package registry

import (
	"bytes"
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
  var got []byte
  var b bytes.Buffer
  ReadMagicBytes(&b, got)

  want := []byte{0x1F, 0x1B}
  assert.Equal(t, want, got)
}
