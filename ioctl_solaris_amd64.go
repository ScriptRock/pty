package pty

import (
	"unsafe"

	"golang.org/x/sys/unix"
)

func init() {
	a := int(123)
	b := uintptr(123)
	if unsafe.Sizeof(a) != unsafe.Sizeof(b) {
		panic("cannot cast between int and uintptr")
	}
}

func ioctl(fd, cmd, ptr uintptr) error {
	return unix.IoctlSetInt(int(fd), uint(cmd), int(ptr))
}
