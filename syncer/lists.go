package syncer

import (
	"os"
	"path/filepath"

	"github.com/pkg/sftp"
)

func containsFile(s []os.FileInfo, e os.FileInfo) bool {
	for _, a := range s {
		if a.Name() == e.Name() {
			return true
		}
	}
	return false
}

func recursiveRemoteList(conn *sftp.Client, root string) []os.FileInfo {
	list := []os.FileInfo{}
	walker := conn.Walk(root)
	for walker.Step() {
		if err := walker.Err(); err != nil {
			panic(err)
		}
		info, err := conn.Stat(walker.Path())
		if err != nil {
			panic(err)
		}
		if walker.Path() == root {
			continue
		}
		list = append(list, info)
	}
	return list
}

func recursiveLocalList(root string) []os.FileInfo {
	list := []os.FileInfo{}
	err := filepath.Walk(root,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if path == root {
				return nil
			}
			list = append(list, info)
			return nil
		})
	if err != nil {
		panic(err)
	}
	return list
}
