package hash

import (
    "crypto/sha1"
)
type Hash_sha1 struct {}

func (o Hash_sha1) GetHash(s string) []byte{
    h := sha1.New()

    h.Write([]byte(s))

    return h.Sum(nil)
}