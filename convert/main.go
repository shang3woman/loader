package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"
	"encoding/hex"
	"bytes"
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
	var buffer bytes.Buffer
	buffer.WriteString(hex.EncodeToString(key))
	buffer.WriteString("\n")
	payload, err := os.ReadFile(os.Args[1])
	if err != nil {
		fmt.Println(err)
		return
	}
	xor(key, payload)
	buffer.WriteString(hex.EncodeToString(payload))
	if err := os.WriteFile("sc.txt", buffer.Bytes(), os.ModePerm); err != nil {
		fmt.Println(err)
	}
}
