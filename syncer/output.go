package syncer

import (
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/vbauerster/mpb"
	"github.com/vbauerster/mpb/decor"
)

type RootSize struct {
	Count int64
	Size  int64
}

type RootStat struct {
	Size       int64
	Downloaded int64
}

func overallSize(files []FileInfo) (size int64) {
	for _, info := range files {
		size += info.Info.Size()
	}
	return
}

func findRoots(files []FileInfo) map[string]*RootSize {
	rootSet := map[string]*RootSize{}
	for _, file := range files {
		root := strings.Split(file.RelativePath, string(filepath.Separator))[0]
		if _, ok := rootSet[root]; !ok {
			rootSet[root] = &RootSize{Count: 1, Size: file.Info.Size()}
		} else {
			rootSet[root].Count++
			rootSet[root].Size += file.Info.Size()
		}
	}
	return rootSet
}

func truncateString(str string, num int) string {
	bnoden := str
	if len(str) > num {
		if num > 3 {
			num -= 3
		}
		bnoden = str[0:num] + "..."
	}
	return bnoden
}

func addBar(p *mpb.Progress, name string, size int64) *mpb.Bar {
	return p.AddBar(size,
		mpb.PrependDecorators(
			decor.CountersKibiByte(truncateString(name, 20)+": % 6.1f / % 6.1f"),
		),
		mpb.AppendDecorators(
			decor.EwmaETA(decor.ET_STYLE_MMSS, 1024),
			decor.Name(" | "),
			decor.AverageSpeed(decor.UnitKiB, "% .2f"),
		),
	)
}

func outputWorker(files []FileInfo, progress <-chan FileProgressInfo, wg *sync.WaitGroup) {
	defer wg.Done()
	size := overallSize(files)
	roots := findRoots(files)

	p := mpb.New()

	start := time.Now()
	totalProgress := addBar(p, "Total", size)

	bars := map[string]*mpb.Bar{}

	rootStats := map[string]*RootStat{}

	for root, stats := range roots {
		bars[root] = addBar(p, root, stats.Size)
		rootStats[root] = &RootStat{Size: stats.Size}
	}

	for p := range progress {
		root := strings.Split(p.Info.RelativePath, string(filepath.Separator))[0]

		totalProgress.IncrBy(int(p.Progress.Increment), time.Since(start))

		bars[root].IncrBy(int(p.Progress.Increment), time.Since(start))
		rootStats[root].Downloaded += p.Progress.Increment
	}
	p.Wait()
}
