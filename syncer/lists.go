package syncer

import (
	"os"
	"path/filepath"
	"sort"

	"github.com/pkg/sftp"
)

type FileInfo struct {
	FullPath     string
	RelativePath string
	Info         os.FileInfo
}

func recursiveRemoteList(conn *sftp.Client, root string) []FileInfo {
	list := []FileInfo{}
	walker := conn.Walk(root)
	for walker.Step() {
		if err := walker.Err(); err != nil {
			panic(err)
		}
		info, err := conn.Stat(walker.Path())
		if os.IsNotExist(err) {
			continue // Can't read file stats, most likely it was deleted, skipping...
		}
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
	sort.Slice(list, func(i, j int) bool {
		return list[i].Info.Size() < list[j].Info.Size()
	})
	return list
}
