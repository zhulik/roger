package syncer

import (
	"sync"

	"github.com/pkg/sftp"
)

func downloadWorker(conn *sftp.Client, jobs <-chan DownloadJob, progress chan<- ProgressInfo, wg *sync.WaitGroup) {
	defer wg.Done()

	for job := range jobs {
		download(conn, job.Info, job.LocalPath, progress)
	}
}
