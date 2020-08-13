package main

import (
    "github.com/kazuminn/dfchange/hash"
    "fmt"
)

func main() {
    s := "sha1 this string"

    sha1 := hash.Hash_sha1{}
    bs := sha1.GetHash(s)


    fmt.Printf("%x\n", bs)
}
