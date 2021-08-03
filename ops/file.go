package ops

import (
	"context"
	"expinc/sunagent/common"
	"expinc/sunagent/log"
	"fmt"
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

func GetFileContent(ctx context.Context, path string) (content []byte, err error) {
	content, err = ioutil.ReadFile(path)
	if nil != err {
		log.ErrorCtx(ctx, err)
	}
	return
}

func WriteFile(ctx context.Context, path string, content []byte, isDir bool, overwrite bool) (meta FileMeta, err error) {
	if false == overwrite {
		// check if file exists
		_, err = os.Stat(path)
		if nil == err {
			err = common.NewError(common.ErrorUnexpected, fmt.Sprintf("%s already exists", path))
			log.ErrorCtx(ctx, err)
			return
		}
		if !os.IsNotExist(err) {
			log.ErrorCtx(ctx, err)
			return
		}
	}

	if false == isDir {
		err = ioutil.WriteFile(path, content, fs.ModePerm)
	} else {
		err = os.MkdirAll(path, fs.ModeDir)
	}

	if nil != err {
		log.ErrorCtx(ctx, err)
	} else {
		metas, err := GetFileMetas(ctx, path, false)
		if err == nil {
			meta = metas[0]
		}
	}
	return
}

func DeleteFile(ctx context.Context, path string, recursive bool) error {
	// Get file info of the path
	pathInfo, err := os.Stat(path)
	if nil != err {
		log.ErrorCtx(ctx, err)
		return err
	}

	// execute removal
	if !pathInfo.IsDir() || !recursive {
		err = os.Remove(path)
	} else {
		err = os.RemoveAll(path)
	}

	return err
}
