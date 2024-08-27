package gpadkm

import (
	"fmt"
	"os"
	"strings"
	"syscall"
	"unsafe"
)

type AvailableEventDevice struct {
	Name string
	Path string
}

func GetAvailableEventDevices() ([]AvailableEventDevice, error) {
	devices := make([]AvailableEventDevice, 0)

	// Open the directory
	dir, err := os.Open("/dev/input")
	if err != nil {
		return nil, err
	}
	defer dir.Close()

	// Read the directory
	files, err := dir.Readdir(-1)
	if err != nil {
		return nil, err
	}

	// Iterate over the files
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		if strings.HasPrefix(file.Name(), "event") {
			path := fmt.Sprintf("/dev/input/%s", file.Name())
			fd, err := syscall.Open(path, syscall.O_RDONLY|syscall.O_NONBLOCK, 0)
			if err != nil {
				continue
			}
			name := make([]byte, 256)
			_, _, errno := syscall.Syscall(syscall.SYS_IOCTL, uintptr(fd), uintptr(syscallEVIOCGNAME(256)), uintptr(unsafe.Pointer(&name[0])))
			if errno != 0 {
				syscall.Close(fd)
				continue
			}
			syscall.Close(fd)
			devices = append(devices, AvailableEventDevice{Name: string(name), Path: path})

		}
	}

	return devices, nil
}
