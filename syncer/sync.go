package syncer

import (
	"log"
	"net/url"
	"os"

	"github.com/pkg/sftp"
)

func Sync(conn *sftp.Client, local string, remote *url.URL, workers int) {
	log.Println("Reading file lists...")
	localFiles := recursiveLocalList(local)
	remoteFiles := recursiveRemoteList(conn, remote.Path)

	filesToSync := []os.FileInfo{}

	for _, info := range remoteFiles {
		if containsFile(localFiles, info) {
			continue
		}
		filesToSync = append(filesToSync, info)
	}
	log.Printf("Found %d files to sync", len(filesToSync))
	log.Printf("Syncing with %d workers...", workers)
}
