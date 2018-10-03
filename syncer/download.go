package syncer

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/pkg/sftp"
)

type ProgressInfo struct {
	Info     FileInfo
	Progress int64
}

func preparePath(fPath string) error {
	parent := filepath.Dir(fPath)
	return os.MkdirAll(parent, os.ModePerm)
}

func progressCopy(r io.Reader, w io.Writer, progress chan<- int64) {
	tee := io.TeeReader(r, w)
	buf := make([]byte, 1<<20)
	var copied int64
	for {
		n, err := tee.Read(buf)
		if err == io.EOF {
			copied += int64(n)
			progress <- copied
			break
		}
		if err != nil {
			panic(err)
		}
		copied += int64(n)
		progress <- copied
	}
	close(progress)
}

func download(conn *sftp.Client, info FileInfo, to string, progress chan<- ProgressInfo) {
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
	pChan := make(chan int64)
	go progressCopy(f, t, pChan)
	for p := range pChan {
		progress <- ProgressInfo{Info: info, Progress: p}
	}
	os.Rename(tmpPath, to)
}
