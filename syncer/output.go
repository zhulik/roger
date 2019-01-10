package syncer

import (
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/vbauerster/mpb"
	"github.com/vbauerster/mpb/decor"
	"github.com/zhulik/roger/storage"
)

func overallSize(files []FileInfo) (size int64) {
	for _, info := range files {
		size += info.Info.Size()
	}
	return
}

func findRoots(files []FileInfo) map[string]int64 {
	rootSet := map[string]int64{}
	for _, file := range files {
		root := strings.Split(file.RelativePath, string(filepath.Separator))[0]
		if _, ok := rootSet[root]; !ok {
			rootSet[root] = file.Info.Size()
		} else {
			rootSet[root] += file.Info.Size()
		}
	}
	return rootSet
}

func truncateString(str string, num int) string {
	bnoden := []rune(str)
	if len(bnoden) < num {
		bnoden = append(bnoden, []rune(strings.Repeat(" ", num-len(bnoden)))...)
	}
	if len(bnoden) > num {
		if num > 3 {
			num -= 3
		}
		bnoden = append(bnoden[0:num], []rune("...")...)
	}
	return string(bnoden)
}

func addBar(p *mpb.Progress, name string, size int64) *mpb.Bar {
	return p.AddBar(size,
		mpb.PrependDecorators(
			decor.CountersKibiByte(truncateString(name, 30)+": %-8.1f / %-8.1f"),
		),
		mpb.AppendDecorators(
			decor.AverageETA(decor.ET_STYLE_MMSS),
			decor.Name(" | "),
			decor.AverageSpeed(decor.UnitKiB, "% .1f"),
		),
	)
}

func updateStorage(files []FileInfo, progress <-chan FileProgressInfo, wg *sync.WaitGroup) {
	defer wg.Done()
	size := overallSize(files)
	roots := findRoots(files)
	storage := storage.GetInstance()
	storage.Initialize(size, roots)

	for p := range progress {
		root := strings.Split(p.Info.RelativePath, string(filepath.Separator))[0]

		storage.Total.Progress += p.Progress
		storage.Files[root].Progress += p.Progress
	}
	storage.Reset()
}

func outputWorker(wg *sync.WaitGroup) {
	defer wg.Done()
	storage := storage.GetInstance()

	p := mpb.New(mpb.WithRefreshRate(1 * time.Second))

	totalProgress := addBar(p, "Total", storage.Total.Total)

	bars := map[string]*mpb.Bar{}

	for root, stats := range storage.Files {
		bars[root] = addBar(p, root, stats.Total)
	}

	var start time.Time

	for storage.Total != nil {
		start = time.Now()

		totalProgress.IncrBy(int(storage.Total.Progress), time.Since(start))

		for root, stats := range storage.Files {
			bars[root].IncrBy(int(stats.Progress), time.Since(start))
		}
		time.Sleep(1 * time.Second)
	}
	p.Wait()
}
