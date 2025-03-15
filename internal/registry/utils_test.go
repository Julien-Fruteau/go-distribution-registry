package registry

import (
	"testing"

	"github.com/opencontainers/go-digest"
	"github.com/stretchr/testify/assert"
)

// digest pseudo code calculation
// let C = 'a small string'
// let B = sha256(C)
// let D = 'sha256:' + EncodeHex(B)
// let ID(C) = D
func TestGetDigestDirectly(t *testing.T) {
	digest := digest.FromBytes([]byte("hello"))
	assert.Equal(t, "sha256:2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824", digest.String())
}

func TestD(t *testing.T) {
var j = `{
 "schemaVersion": 2,
 "mediaType": "application/vnd.oci.image.manifest.v1+json",

	"config": {
	  "mediaType": "application/vnd.oci.image.config.v1+json",
	  "digest": "sha256:1649f157365545ac4b8ec167619fb18d2b61f802776e39e46a8156f39762615e",
	  "size": 11148
	},

 "layers": [

	{
	  "mediaType": "application/vnd.oci.image.layer.v1.tar+gzip",
	  "digest": "sha256:9d1c7dcd50f5547c998ed553485c4c8ef1bcba72abb1b70c4f7de74572c54278",
	  "size": 145483495
	},
	{
	  "mediaType": "application/vnd.oci.image.layer.v1.tar+gzip",
	  "digest": "sha256:b9be66bfe7f92b5c42a47c6353d0dfb1f7b9610a9479752228d5f1fe00c100fc",
	  "size": 2094433
	},
	{
	  "mediaType": "application/vnd.oci.image.layer.v1.tar+gzip",
	  "digest": "sha256:08d8d343d6a4c6fb7033d42667d66a88368e95b0f1ee288621dbaf24149d33ca",
	  "size": 178
	},
	{
	  "mediaType": "application/vnd.oci.image.layer.v1.tar+gzip",
	  "digest": "sha256:e6c0e3d5828e19ef46e585f10e2af75e11be87f42432301fb92df598d2d2d092",
	  "size": 477195673
	},
	{
	  "mediaType": "application/vnd.oci.image.layer.v1.tar+gzip",
	  "digest": "sha256:e8ff69f6858575d6e0a8be832b30716e41bce7379acee914fa91da26533a484a",
	  "size": 477197575
	},
	{
	  "mediaType": "application/vnd.oci.image.layer.v1.tar+gzip",
	  "digest": "sha256:107aba61455803961e2bf3981ee15312ff09177dfa623fe785f4759b63afa9a5",
	  "size": 7267273
	},
	{
	  "mediaType": "application/vnd.oci.image.layer.v1.tar+gzip",
	  "digest": "sha256:4f4fb700ef54461cfa02571ae0db9a0dc1e0cdb5577484a6d75e68dc38e8acc1",
	  "size": 32
	}
`
  d := digest.FromBytes([]byte(j))
  assert.Equal(t, "sha256:9d1c7dcd50f5547c998ed553485c4c8ef1bcba72abb1b70c4f7de74572c54278", d.String())
}
