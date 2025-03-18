# go-distribution-registry

A distribution registry client library

Target :

- [x] collect all repositories from the distribution registry

- [ ] inspect on existing tag can lead to (all) layers 404
- [ ] not here remove all cluster images tag from images tag to delete
- [ ] wip: for all remaining tag, inspect and extract created date
- [ ] remove from the list the N more recent tags

- [ ] wip: for remaining tags, call distribution delete tag

- [ ] then a registry garbage collect should be called/executed
  NB: since current v2 is buggy, consider doing it manually
  while v3 is officially released

 🔥: not direct http call to get the list of all available blobs

  in classic storage :

    1- Recursively scans the blobs directory.
    2- Extracts the digest from the filesystem path.
    3- Prints all available blob SHA256 digests.

```go
  package main

  import (
    "fmt"
    "os"
    "path/filepath"
    "strings"
  )

  const blobDir = "/var/lib/registry/docker/registry/v2/blobs/sha256"

  func main() {
    err := filepath.Walk(blobDir, func(path string, info os.FileInfo, err error) error {
      if err != nil {
        return err
      }
      if !info.IsDir() && len(path) > len(blobDir)+4 { // Ignore directories, expect sha256/<2-char>/<rest>
        relPath, _ := filepath.Rel(blobDir, path)
        parts := strings.Split(relPath, string(os.PathSeparator))
        if len(parts) == 2 { // Should be of the form sha256/<first2>/<remaining_digest>
          fmt.Printf("sha256:%s%s\n", parts[0], parts[1])
        }
      }
      return nil
    })

    if err != nil {
      fmt.Println("Error:", err)
    }
  }
```

blobs can be either a json describing the container config, the list of manifest or tar+gzip file actually containing
the image container layer

all blobs are stored in the same location, so to identify tar+gzip blobs an additional identification step is required

Every file type has magic bytes at the beginning.

in bash :

```sh
xxd -l 4 myfile.tar.gz
# 00000000: 1f8b 08xx
# 1F 8B → Gzip signature
# 08    → Compression method (Deflate)
```

```go
package main

import (
 "bytes"
 "compress/gzip"
 "fmt"
 "os"
)

func isGzipFile(filename string) bool {
 file, err := os.Open(filename)
 if err != nil {
  fmt.Println("Error opening file:", err)
  return false
 }
 defer file.Close()

 // Read first few bytes 2, or 4?
 buf := make([]byte, 2)
 _, err = file.Read(buf)
 if err != nil {
  return false
 }

 // Check for gzip magic number (1F 8B)
 return bytes.Equal(buf, []byte{0x1F, 0x8B})
}

```
