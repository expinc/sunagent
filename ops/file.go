package ops

import (
	"context"
	"expinc/sunagent/common"
	"expinc/sunagent/log"
	"io/fs"
	"io/ioutil"
	"os"
	"time"
)

type FileMeta struct {
	Name             string    `json:"name"`
	Size             int64     `json:"size"`
	LastModifiedTime time.Time `json:"lastModifiedTime"`
	Owner            string    `json:"owner"`
	Mode             string    `json:"mode"`
}

func fileInfoToMeta(ctx context.Context, info fs.FileInfo) (meta FileMeta) {
	meta.Name = info.Name()
	meta.Size = info.Size()
	meta.LastModifiedTime = info.ModTime()
	meta.Mode = info.Mode().String()
	meta.Owner = getOwner(ctx, info)
	return
}

func GetFileMetas(ctx context.Context, path string, listIfDir bool) (metas []FileMeta, err error) {
	metas = make([]FileMeta, 0)

	// Get file info of the path
	var pathInfo os.FileInfo
	pathInfo, err = os.Stat(path)
	if nil != err {
		log.ErrorCtx(ctx, err)
		err = common.NewError(common.ErrorNotFound, err.Error())
		return
	}

	// Get info of contained files if the path is directory and need to list it
	var infos []fs.FileInfo
	if pathInfo.IsDir() && listIfDir {
		infos, err = ioutil.ReadDir(path)
	} else {
		infos = []fs.FileInfo{pathInfo}
	}
	if nil != err {
		log.ErrorCtx(ctx, err)
		return
	}

	// Get meta of listed files
	for _, info := range infos {
		metas = append(metas, fileInfoToMeta(ctx, info))
	}

	return
}
