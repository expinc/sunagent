package ops

import (
	"context"
	"expinc/sunagent/log"
	"os"

	"golang.org/x/sys/windows"
)

func getOwner(ctx context.Context, info os.FileInfo) string {
	errorReturn := func(err error) string {
		log.ErrorCtx(ctx, err)
		return ""
	}

	namep, err := windows.UTF16PtrFromString(info.Name())
	if nil != err {
		return errorReturn(err)
	}

	handler, err := windows.CreateFile(namep, windows.GENERIC_READ, windows.FILE_SHARE_READ, nil, windows.OPEN_EXISTING, windows.FILE_FLAG_BACKUP_SEMANTICS, 0)
	if nil != err {
		return errorReturn(err)
	}
	defer windows.CloseHandle(handler)

	secInfo, err := windows.GetSecurityInfo(handler, windows.SE_FILE_OBJECT, windows.OWNER_SECURITY_INFORMATION)
	if nil != err {
		return errorReturn(err)
	}

	ownerSid, _, err := secInfo.Owner()
	if nil != err {
		return errorReturn(err)
	}

	account, _, _, err := ownerSid.LookupAccount("")
	if nil != err {
		return errorReturn(err)
	}
	return account
}
