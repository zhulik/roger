package syncer

import (
	"log"
	"path/filepath"

	"github.com/pkg/sftp"
)

func Sync(conn *sftp.Client, local string, remote string, workers int) {
	log.Println("Reading file lists...")
	localFiles := recursiveLocalList(local)
	remoteFiles := recursiveRemoteList(conn, remote)

	log.Println("Comparing...")

	filesToSync := []FileInfo{}

	for _, info := range remoteFiles {
		if containsFile(localFiles, info) {
			continue
		}
		filesToSync = append(filesToSync, info)
	}
	log.Printf("Found %d files to sync", len(filesToSync))
	// log.Printf("Syncing with %d workers...", workers)
	log.Println("Naive syncing...")
	for _, info := range filesToSync {
		localPath := filepath.Join(local, info.RelativePath)
		log.Printf("Syncing %s to %s", info.FullPath, localPath)
		download(conn, info.FullPath, localPath)
	}

}
