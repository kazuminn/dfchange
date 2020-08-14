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
var root = "/home/kazumi/go/src/github.com/kazuminn/dfchange"

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

func exists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

func main() {

	ok := exists("./result.json")
	if ok {
        //current
		err := filepath.Walk(root, scan)
		if err != nil {
			fmt.Println(err)
        }
		currentConfig := configs
        currentMap := map[string]string{}
        for _, r := range currentConfig {
            currentMap[r.Path] = string(r.Hash)
		}

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
        prevMap := map[string]string{}
        for _, r := range prevConfig {
            prevMap[r.Path] = string(r.Hash)
		}

        //detect
        for path, hash := range prevMap {
            _, ok := currentMap[path]
            if(currentMap[path] != hash) {
			    fmt.Println(path)
            }else if(!ok){
			    fmt.Println(path)
            }
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
