package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/kazuminn/dfchange/hash"
)

type config struct {
	Path string `json:"path"`
	Hash []byte `json:"hash"`
}

var configs = []config{}

func scan(path string, info os.FileInfo, err error) error {
	if info.IsDir() {
		// 特定のディレクトリ以下を無視する場合は
		// return filepath.SkipDir
		return nil
	}

	f, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
	}

	// 一気に全部読み取り
	b, err := ioutil.ReadAll(f)

	sha1 := hash.Hash_sha1{}
	bs := sha1.GetHash(string(b))
	configs = append(configs, config{path, bs})
	f.Close()

	return nil
}

func main() {

	root := "/home/kazumi/go/src/github.com/kazuminn/dfchange"

	err := filepath.Walk(root, scan)

	if err != nil {
		fmt.Println(err)
	}

	s, _ := json.Marshal(configs)
	fmt.Println(string(s))
}
