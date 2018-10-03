package syncer

import (
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/pkg/sftp"
)

type DownloadJob struct {
	Info      FileInfo
	LocalPath string
}

func Sync(conn *sftp.Client, local string, remote string, workers int) {
	log.Println("Reading file list...")
	remoteFiles := recursiveRemoteList(conn, remote)

	log.Println("Comparing...")

	filesToSync := []FileInfo{}

	for _, info := range remoteFiles {
		if _, err := os.Stat(filepath.Join(local, info.RelativePath)); os.IsNotExist(err) {
			filesToSync = append(filesToSync, info)
		}
	}
	log.Printf("Found %d files to sync", len(filesToSync))
	// log.Printf("Syncing with %d workers...", workers)
	log.Println("Naive syncing...")
	jobs := make(chan DownloadJob)
	progress := make(chan ProgressInfo)
	wg := sync.WaitGroup{}
	wg.Add(1)
	go downloadWorker(conn, jobs, progress, &wg)

	go func() {
		for _, info := range filesToSync {
			localPath := filepath.Join(local, info.RelativePath)
			// log.Printf("Syncing %s to %s", info.FullPath, localPath)
			jobs <- DownloadJob{Info: info, LocalPath: localPath}
		}
		log.Println("Closing jobs...")
		close(jobs)
	}()

	logWg := sync.WaitGroup{}
	logWg.Add(1)
	go outputWorker(filesToSync, progress, &logWg)

	log.Println("Waiting for jobs..")
	wg.Wait()

	log.Println("Closing progress...")
	close(progress)

	log.Println("Waiting for progress...")
	logWg.Wait()
}
