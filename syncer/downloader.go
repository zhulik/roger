package syncer

import (
	"io"
	"os"
	"path/filepath"

	"github.com/pkg/sftp"
)

func preparePath(fPath string) error {
	parent := filepath.Dir(fPath)
	return os.MkdirAll(parent, os.ModePerm)
}

func download(conn *sftp.Client, from, to string) {
	f, err := conn.Open(from)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	err = preparePath(to)
	if err != nil {
		panic(err)
	}
	t, err := os.Create(to)
	if err != nil {
		panic(err)
	}
	defer t.Close()
	_, err = io.Copy(t, f)
	if err != nil {
		panic(err)
	}
}
