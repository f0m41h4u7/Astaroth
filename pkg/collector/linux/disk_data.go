package linux

import (
	"errors"
	"os/exec"
	"strconv"
	"strings"

	"github.com/f0m41h4u7/Astaroth/pkg/api"
)

var (
	ErrWrongDiskData    = errors.New("cannot parse disk data")
	ErrCannotParseUsed  = errors.New("cannot parse used percentage")
	ErrCannotParseInode = errors.New("cannot parse inode percentage")
)

func GetDiskData() (*api.DiskData, error) {
	mb, err := exec.Command("df", "-k").Output()
	if err != nil {
		return nil, err
	}
	inode, err := exec.Command("df", "-i").Output()
	if err != nil {
		return nil, err
	}
	return parseDiskData(string(mb), string(inode))
}

func parseDiskData(mb string, inode string) (*api.DiskData, error) {
	mbLines := strings.Split(mb, "\n")[1:]
	inodeLines := strings.Split(inode, "\n")[1:]
	disk := api.DiskData{
		Data: make([]*api.FilesystemData, len(mbLines)),
	}

	if (len(mb) == 0) || (len(inode) == 0) {
		return &disk, ErrWrongDiskData
	}

	if len(mbLines) != len(inodeLines) {
		return &disk, ErrWrongDiskData
	}
	for i := 0; i < len(mbLines); i++ {
		fields := strings.Fields(strings.TrimSpace(mbLines[i]))
		fs := fields[0]
		used, err := strconv.ParseInt(strings.Split(fields[4], "%")[0], 10, 64)
		if err != nil {
			return &disk, ErrCannotParseUsed
		}
		fields = strings.Fields(strings.TrimSpace(inodeLines[i]))
		if fs != fields[0] {
			return &disk, ErrWrongDiskData
		}
		iused, err := strconv.ParseInt(strings.Split(fields[4], "%")[0], 10, 64)
		if err != nil {
			return &disk, ErrCannotParseInode
		}

		disk.Data[i] = &api.FilesystemData{
			Filesystem: fs,
			Used:       used,
			Inode:      iused,
		}
	}
	return &disk, nil
}
