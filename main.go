package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"syscall"
	"strconv"

	"github.com/kazuminn/dfchange/hash"
)

type config struct {
	Path string `json:"path"`
	Hash []byte `json:"hash"`
	Permission string `json:"permission"`
	Uid string `json:"uid"`
	Gid string `json:"gid"`
}

var currentConfig []config
var configs = []config{}
var prevMap = map[string]string{}
var root = "/"

func scan(path string, info os.FileInfo, err error) error {
	if path == "/swapfile" {
		return nil
	}
	if info.IsDir() {
		if path == "/var/log" || path == "/sys/kernel" || path == "/snap" || path == "/dev" || path == "~/components" || path == "/proc" || path == "/run" {
			return filepath.SkipDir
		}
	}

	var uid string
	var gid string
	if stat, ok := info.Sys().(*syscall.Stat_t); ok {
		uid = strconv.Itoa(int(stat.Uid))
		gid = strconv.Itoa(int(stat.Gid))
	}

	f, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(path)


	// 一気に全部読み取り
	b, err := ioutil.ReadAll(f)

	sha1 := hash.Hash_sha1{}
	bs := sha1.GetHash(string(b))
	configs = append(configs, config{path, bs, info.Mode().String(), uid, gid})
	f.Close()

	return nil
}

func detect(path string, info os.FileInfo, err error) error {
	if path == "/swapfile" || path == "/sys/power/wakeup_count" {
		return nil
	}
	if info.IsDir() {
		if path == "/sys/devices" || path == "/sys/fs" || path == "/var/log" || path == "/sys/kernel" || path == "/snap" || path == "/dev" || path == "~/components" || path == "/proc" || path == "/run" {
			return filepath.SkipDir
		}
	}
	f, _ := os.Open(path)


	// 一気に全部読み取り
	b, err := ioutil.ReadAll(f)

	var uid string
	var gid string
	if stat, ok := info.Sys().(*syscall.Stat_t); ok {
		uid = strconv.Itoa(int(stat.Uid))
		gid = strconv.Itoa(int(stat.Gid))
	}

	sha1 := hash.Hash_sha1{}
	bs := sha1.GetHash(string(b))
	value, ok := prevMap[path]
	if ok {
		if value != string(bs) {
			fmt.Println("changed file :%s", path)
		}
	} else { //detect new file
		fmt.Println("changed file :%s", path)
	}
	configs = append(configs, config{path, bs, info.Mode().String(), uid, gid})
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

		//detect
		err = filepath.Walk(root, detect)
		if err != nil {
			fmt.Println(err)
		}

		//file write
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
