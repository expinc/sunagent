package ops

import (
	"expinc/sunagent/log"
	"io/fs"
	"io/ioutil"
	"os"
	"time"
)

type FileMeta struct {
	Name             string    `json:"version"`
	Size             int64     `json:"size"`
	LastModifiedTime time.Time `json:"lastModifiedTime"`
	Owner            string    `json:"owner"`
	Mode             string    `json:"mode"`
}

func fileInfoToMeta(info fs.FileInfo) (meta FileMeta) {
	meta.Name = info.Name()
	meta.Size = info.Size()
	meta.LastModifiedTime = info.ModTime()
	meta.Mode = info.Mode().String()
	meta.Owner = getOwner(info)
	return
}

func GetFileMetas(path string, listIfDir bool) (metas []FileMeta, err error) {
	metas = make([]FileMeta, 0)

	// Get file info of the path
	var pathInfo os.FileInfo
	pathInfo, err = os.Stat(path)
	if nil != err {
		log.Error(err)
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
		log.Error(err)
		return
	}

	// Get meta of listed files
	for _, info := range infos {
		metas = append(metas, fileInfoToMeta(info))
	}

	return
}
