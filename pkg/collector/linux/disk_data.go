package linux

import (
	"errors"
	"os/exec"
	"strconv"
	"strings"

	"github.com/f0m41h4u7/Astaroth/pkg/api"
)

var (
	errWrongDiskData    = errors.New("cannot parse disk data")
	errCannotParseUsed  = errors.New("cannot parse used percentage")
	errCannotParseInode = errors.New("cannot parse inode percentage")
)

// GetDiskData collects information about disk usage.
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
		return &disk, errWrongDiskData
	}

	if len(mbLines) != len(inodeLines) {
		return &disk, errWrongDiskData
	}
	for i := 0; i < len(mbLines); i++ {
		fields := strings.Fields(strings.TrimSpace(mbLines[i]))
		fs := fields[0]
		used, err := strconv.ParseInt(strings.Split(fields[4], "%")[0], 10, 64)
		if err != nil {
			return &disk, errCannotParseUsed
		}
		fields = strings.Fields(strings.TrimSpace(inodeLines[i]))
		if fs != fields[0] {
			return &disk, errWrongDiskData
		}
		iused, err := strconv.ParseInt(strings.Split(fields[4], "%")[0], 10, 64)
		if err != nil {
			return &disk, errCannotParseInode
		}

		disk.Data[i] = &api.FilesystemData{
			Filesystem: fs,
			Used:       used,
			Inode:      iused,
		}
	}

	return &disk, nil
}
