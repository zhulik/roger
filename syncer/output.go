package syncer

import (
	"log"
	"sync"
)

func overallSize(files []FileInfo) (size int64) {
	for _, info := range files {
		size += info.Info.Size()
	}
	return
}

func outputWorker(files []FileInfo, progress <-chan ProgressInfo, wg *sync.WaitGroup) {
	defer wg.Done()
	count := len(files)
	size := overallSize(files)

	var downloadedCount int64
	var downloadedBytes int64

	for p := range progress {
		if p.Progress == p.Info.Info.Size() {
			downloadedCount++
		}
		downloadedBytes += p.Progress
		log.Printf("Progress count=%d/%d bytes=%d/%d percent=%f", downloadedCount, count, downloadedBytes, size, float64(downloadedBytes)/float64(size)*100)
	}
}
