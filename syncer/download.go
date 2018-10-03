package syncer

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/pkg/sftp"
)

type FileProgressInfo struct {
	Info     FileInfo
	Progress ProgressInfo
}

type ProgressInfo struct {
	Increment int64
	Total     int64
}

func preparePath(fPath string) error {
	parent := filepath.Dir(fPath)
	return os.MkdirAll(parent, os.ModePerm)
}

func progressCopy(r io.Reader, w io.Writer, progress chan<- ProgressInfo) {
	tee := io.TeeReader(r, w)
	buf := make([]byte, 1<<20)
	var copied int64
	for {
		n, err := tee.Read(buf)
		if err == io.EOF {
			copied += int64(n)
			progress <- ProgressInfo{Increment: int64(n), Total: copied}
			break
		}
		if err != nil {
			panic(err)
		}
		copied += int64(n)
		progress <- ProgressInfo{Increment: int64(n), Total: copied}
	}
	close(progress)
}

func download(conn *sftp.Client, info FileInfo, to string, progress chan<- FileProgressInfo) {
	f, err := conn.Open(info.FullPath)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	err = preparePath(to)
	if err != nil {
		panic(err)
	}
	tmpPath := fmt.Sprintf("%s.crdownload", to)
	t, err := os.Create(tmpPath)
	if err != nil {
		panic(err)
	}
	defer t.Close()
	pChan := make(chan ProgressInfo)
	go progressCopy(f, t, pChan)
	for p := range pChan {
		progress <- FileProgressInfo{Info: info, Progress: p}
	}
	os.Rename(tmpPath, to)
}
