package syncer

import (
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

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
	if len(filesToSync) == 0 {
		log.Println("Nothing to sync, exiting...")
		return
	}
	log.Printf("Found %d files to sync", len(filesToSync))
	log.Printf("Syncing with %d workers...", workers)
	jobs := make(chan DownloadJob)
	progress := make(chan FileProgressInfo)
	wg := sync.WaitGroup{}
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go downloadWorker(conn, jobs, progress, &wg)
	}

	go func() {
		for _, info := range filesToSync {
			localPath := filepath.Join(local, info.RelativePath)
			jobs <- DownloadJob{Info: info, LocalPath: localPath}
		}
		close(jobs)
	}()

	logWg := sync.WaitGroup{}
	logWg.Add(2)
	go updateStorage(filesToSync, progress, &logWg)
	time.Sleep(1 * time.Second)
	go outputWorker(&logWg)

	wg.Wait()
	close(progress)
	logWg.Wait()
}
