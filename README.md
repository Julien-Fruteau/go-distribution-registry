# go-distribution-registry

A distribution registry client library

Target : clean up repository tag and image layers

## TODO

### Http Client

- [x] collect all repositories from the distribution registry
- [x] for tag, inspect and extract created date
  - [ ] WIP : validate cmd tagsDate
- [ ] remove from the list the N more recent tags

- [x] delete tag, (blobs) layer

~- [ ] then a registry garbage collect should be called/executed
  NB: since current v2 is buggy, consider doing it manually
  while v3 is officially released~

🔥: not direct http call to get the list of all available blobs

### File system Client

- [x] parse filesystem to get all available gzip blobs
  in classic storage :

      1- Recursively scans the blobs directory.
      2- Extracts the digest from the filesystem path.
      3- Prints all available blob SHA256 digests.

```go
 // Approach 3: Using a single []byte (Most Efficient)
 digestsAsSingleSlice := make([]byte, numDigests*32)
 for i := 0; i < numDigests; i++ {
  copy(digestsAsSingleSlice[i*32:(i+1)*32], hash[:])
 }
 fmt.Printf("Memory for single []byte: ~%d MB\n", (len(digestsAsSingleSlice)+int(unsafe.Sizeof(digestsAsSingleSlice)))/(1024*1024))

 // Accessing a digest in the single []byte storage
 index := 10
 start := index * 32
 end := start + 32
 fmt.Printf("Digest at index %d: %x\n", index, digestsAsSingleSlice[start:end])
```
