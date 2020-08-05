package linux

import (
	"errors"
	"os/exec"
	"strconv"
	"strings"
	"sync"

	"github.com/f0m41h4u7/Astaroth/pkg/api"
)

var errWrongDiskData = errors.New("cannot parse disk data")

// GetDiskData collects information about disk usage.
func (c *Collector) getDiskData(wg *sync.WaitGroup, ss *Snapshot) error {
	defer wg.Done()

	mb, err := exec.Command("df", "-k").Output()
	if err != nil {
		return err
	}
	inode, err := exec.Command("df", "-i").Output()
	if err != nil {
		return err
	}
	ss.DiskData, err = parseDiskData(string(mb), string(inode))

	return err
}

func parseDiskData(mb string, inode string) (*api.DiskData, error) {
	mbLines := strings.Split(mb, "\n")[1:]
	inodeLines := strings.Split(inode, "\n")[1:]
	disk := api.DiskData{
		Data: []*api.FilesystemData{},
	}

	if (len(mb) == 0) || (len(inode) == 0) {
		return &disk, errWrongDiskData
	}

	if len(mbLines) != len(inodeLines) {
		return &disk, errWrongDiskData
	}

	for i := 0; i < len(mbLines); i++ {
		fields := strings.Fields(strings.TrimSpace(mbLines[i]))
		if len(fields) == 0 {
			continue
		}
		fs := fields[0]
		used, err := strconv.ParseInt(strings.Split(fields[4], "%")[0], 10, 64)
		if err != nil {
			used = 0
		}
		fields = strings.Fields(strings.TrimSpace(inodeLines[i]))
		if fs != fields[0] {
			return &disk, errWrongDiskData
		}
		iused, err := strconv.ParseInt(strings.Split(fields[4], "%")[0], 10, 64)
		if err != nil {
			iused = 0
		}

		disk.Data = append(disk.Data, &api.FilesystemData{
			Filesystem: fs,
			Used:       used,
			Inode:      iused,
		})
	}

	return &disk, nil
}
