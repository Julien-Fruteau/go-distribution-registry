# distribution-ctl

Interact with distribution registry over http to perform
images tag clean up considering kubernetes cluster images on use

Target :

- [x] parse kubernetes cluster images
- [x] collect all repositories from the distribution registry

- [ ] remove all cluster images tag from images tag to delete
- [ ] wip: for all remaining tag, inspect and extract created date
- [ ] remove from the list the N more recent tags

- [ ] wip: for remaining tags, call distribution delete tag

- [ ] split go-distriubution-registry and go-kubernetes-client as 2 libs and make an app from both lib targeting purpose

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
