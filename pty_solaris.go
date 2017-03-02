package pty

import (
	"fmt"
	"os"
	"os/user"
	"strconv"
	"syscall"
	"unsafe"

	"golang.org/x/sys/unix"
)

const (
	NODEV    = -1
	O_MAXMIN = 0xff // minor device number maximum

	STR   = ('S' << 8) // for strioctl ioctls
	I_STR = (STR | 010)

	ISPTM    = (('P' << 8) | 1) // query for master
	UNLKPT   = (('P' << 8) | 2) // unlock master/slave pair
	PTSSTTY  = (('P' << 8) | 3) // set tty flag
	ZONEPT   = (('P' << 8) | 4) // set zone of master/slave pair
	OWNERPT  = (('P' << 8) | 5) // set owner/group for slave device
	_PTSSTTY = (('P' << 8) | 6) // set tty flag as passed by ldterm
)

func open() (pty, tty *os.File, err error) {
	p, err := os.OpenFile("/dev/ptmx", os.O_RDWR, 0)
	if err != nil {
		return nil, nil, err
	}

	sname, err := ptsname(p)
	if err != nil {
		return nil, nil, err
	}

	err = grantpt(p)
	if err != nil {
		return nil, nil, err
	}

	err = unlockpt(p)
	if err != nil {
		return nil, nil, err
	}

	t, err := os.OpenFile(sname, os.O_RDWR, 0)
	if err != nil {
		return nil, nil, err
	}

	return p, t, nil
}

func ptsdev(f *os.File) (int, error) {
	ioctlData := strioctl{
		ic_cmd:    ISPTM,
		ic_timout: 0,
		ic_len:    0,
		ic_dp:     0,
	}

	if err := ioctl(f.Fd(), I_STR, uintptr(unsafe.Pointer(&ioctlData))); err != nil {
		return 0, err
	}
	if stat, err := f.Stat(); err != nil {
		return 0, err
	} else if sys := stat.Sys(); sys == nil {
		return 0, fmt.Errorf("stat.Sys() is nil")
	} else if stat_t, ok := sys.(*syscall.Stat_t); !ok {
		return 0, fmt.Errorf("stat.Sys() is type %T, %v", stat_t, ok)
	} else {
		return int(uint(stat_t.Rdev) & O_MAXMIN), nil
	}
}

func ptsname(f *os.File) (string, error) {

	dev, err := ptsdev(f)
	if err != nil {
		return "", err
	} else if dev == NODEV {
		return "", fmt.Errorf("no pts device")
	}

	ptsn := fmt.Sprintf("/dev/pts/%d", dev)
	if err := unix.Access(ptsn, 0); err != nil {
		return "", err
	}
	return ptsn, nil
}

func grantpt(f *os.File) error {
	// TODO FIXME: 2nd should be "tty" group number if available
	uid := os.Getuid()
	gid := os.Getgid()
	if ttyGroup, err := user.LookupGroup("tty"); err == nil {
		if n, err := strconv.Atoi(ttyGroup.Gid); err == nil {
			gid = n
		}
	}

	ownerData := pt_own{
		pto_ruid: uid_t(uid),
		pto_rgid: gid_t(gid),
	}
	ioctlData := strioctl{
		ic_cmd:    OWNERPT,
		ic_timout: 0,
		ic_len:    int32(unsafe.Sizeof(ownerData)),
		ic_dp:     uintptr(unsafe.Pointer(&ownerData)),
	}

	if err := ioctl(f.Fd(), I_STR, uintptr(unsafe.Pointer(&ioctlData))); err != nil {
		return fmt.Errorf("grantpt error: %v", err)
	}
	return nil
}

func unlockpt(f *os.File) error {
	ioctlData := strioctl{
		ic_cmd:    UNLKPT,
		ic_timout: 0,
		ic_len:    0,
		ic_dp:     0,
	}
	if err := ioctl(f.Fd(), I_STR, uintptr(unsafe.Pointer(&ioctlData))); err != nil {
		return fmt.Errorf("unlockpt error: %v", err)
	}
	return nil
}
