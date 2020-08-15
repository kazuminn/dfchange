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

var currentConfig []config
var configs = []config{}
var prevMap = map[string]string{}
var root = "/etc"

func scan(path string, info os.FileInfo, err error) error {
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

func detect(path string, info os.FileInfo, err error) error {
	f, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
	}

	// 一気に全部読み取り
	b, err := ioutil.ReadAll(f)

	sha1 := hash.Hash_sha1{}
	bs := sha1.GetHash(string(b))
	value, ok := prevMap[path]
	if ok {
		if value != string(bs) {
			fmt.Println(path)
		}
	} else {
		fmt.Println(path)
	}
	currentConfig = append(currentConfig, config{path, bs})
	f.Close()

	return nil
}

func exists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

func main() {

	ok := exists("./result.json")
	if ok {
		//previous
		b, err := ioutil.ReadFile("./result.json")
		if err != nil {
			fmt.Println(os.Stderr, err)
			os.Exit(1)
		}
		var prevConfig []config
		err = json.Unmarshal(b, &prevConfig)
		if err != nil {
			fmt.Println(os.Stderr, err)
			os.Exit(1)
		}
		for _, r := range prevConfig {
			prevMap[r.Path] = string(r.Hash)
		}

		//current
		err = filepath.Walk(root, detect)
		if err != nil {
			fmt.Println(err)
		}

		//write
		s, _ := json.Marshal(currentConfig)
		err = ioutil.WriteFile("./result.json", s, 0666)
		if err != nil {
			fmt.Println(os.Stderr, err)
			os.Exit(1)
		}

	} else {

		err := filepath.Walk(root, scan)

		if err != nil {
			fmt.Println(err)
		}

		s, _ := json.Marshal(configs)

		err = ioutil.WriteFile("./result.json", s, 0666)
		if err != nil {
			fmt.Println(os.Stderr, err)
			os.Exit(1)
		}
	}
}
