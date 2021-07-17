package ops

import (
	"os"
	"os/user"
	"strconv"
	"syscall"
)

func getOwner(info os.FileInfo) string {
	result := ""
	uid := uint64(info.Sys().(*syscall.Stat_t).Uid)
	usr, _ := user.LookupId(strconv.FormatUint(uid, 10))
	result = usr.Username
	return result
}
