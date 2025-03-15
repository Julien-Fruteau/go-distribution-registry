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

- [ ] then a registry garbage collect should be called/executed
  NB: since current v2 is buggy, consider doing it manually
  while v3 is officially released

- [ ] split go-distriubution-registry and go-kubernetes-client as 2 libs and make an app from both lib targeting purpose
