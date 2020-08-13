package main

import (
	"fmt"
	"github.com/kazuminn/dfchange/hash"
)

type config struct {
	Path string
	Hash []byte
}

func main() {
	s := "sha1 this string"

	sha1 := hash.Hash_sha1{}
	bs := sha1.GetHash(s)

	fmt.Printf("%x\n", bs)
}
