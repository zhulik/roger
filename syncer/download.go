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
	Progress int64
}

func preparePath(fPath string) error {
	parent := filepath.Dir(fPath)
	return os.MkdirAll(parent, os.ModePerm)
}

func progressCopy(r io.Reader, w io.Writer, progress chan<- int64) {
	counter := &WriteCounter{Writer: w, Progress: progress}
	if _, err := io.Copy(counter, r); err != nil {
		panic(err)
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
	pChan := make(chan int64)
	go progressCopy(f, t, pChan)
	for p := range pChan {
		progress <- FileProgressInfo{Info: info, Progress: p}
	}
	if err = os.Rename(tmpPath, to); err != nil {
		panic(nil)
	}
}
