package main

import (
	"fmt"
	"io/ioutil"
	"os"
    "path/filepath"
    "encoding/json"

	"github.com/kazuminn/dfchange/hash"
)

type config struct {
	Path string `json:"path"`
	Hash []byte         `json:"hash"`
}

func main() {

    root := "/home/kazumi/go/src/github.com/kazuminn/dfchange"
    configs := []config{}

	err := filepath.Walk(root,
		func(path string, info os.FileInfo, err error) error {
			if info.IsDir() {
				// 特定のディレクトリ以下を無視する場合は
				// return filepath.SkipDir
				return nil
			}

			f, err := os.Open(path)
			if err != nil {
				fmt.Println("can't read file")
			}

			// 一気に全部読み取り
			b, err := ioutil.ReadAll(f)

			sha1 := hash.Hash_sha1{}
			bs := sha1.GetHash(string(b))
			configs = append(configs, config{ path, bs, })
			f.Close()

			return nil
		})

	if err != nil {
		fmt.Println(1, err)
    }
    
    s, _ := json.Marshal(configs)
			fmt.Println(string(s))
}
