package syncer

import (
	"os"
	"path/filepath"

	"github.com/pkg/sftp"
)

type FileInfo struct {
	FullPath     string
	RelativePath string
	Info         os.FileInfo
}

func containsFile(s []FileInfo, e FileInfo) bool {
	for _, a := range s {
		if a.RelativePath == e.RelativePath {
			return true
		}
	}
	return false
}

func recursiveRemoteList(conn *sftp.Client, root string) []FileInfo {
	list := []FileInfo{}
	walker := conn.Walk(root)
	for walker.Step() {
		if err := walker.Err(); err != nil {
			panic(err)
		}
		info, err := conn.Stat(walker.Path())
		if err != nil {
			panic(err)
		}
		if info.IsDir() {
			continue
		}
		if walker.Path() == root {
			continue
		}
		rPath, err := filepath.Rel(root, walker.Path())
		if err != nil {
			panic(err)
		}
		list = append(list, FileInfo{FullPath: walker.Path(), RelativePath: rPath, Info: info})
	}
	return list
}

func recursiveLocalList(root string) []FileInfo {
	list := []FileInfo{}
	err := filepath.Walk(root,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if path == root {
				return nil
			}
			fullPath := filepath.Join(path, info.Name())
			rPath, err := filepath.Rel(root, fullPath)
			if err != nil {
				panic(err)
			}
			list = append(list, FileInfo{FullPath: fullPath, RelativePath: rPath, Info: info})
			return nil
		})
	if err != nil {
		panic(err)
	}
	return list
}
