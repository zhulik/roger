package syncer

import (
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/vbauerster/mpb"
	"github.com/vbauerster/mpb/decor"
)

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

func outputWorker(files []FileInfo, progress <-chan FileProgressInfo, wg *sync.WaitGroup) {
	defer wg.Done()
	size := overallSize(files)
	roots := findRoots(files)

	p := mpb.New(mpb.WithRefreshRate(1 * time.Second))

	start := time.Now()
	totalProgress := addBar(p, "Total", size)

	bars := map[string]*mpb.Bar{}

	rootStats := map[string]*RootStat{}

	for root, stats := range roots {
		bars[root] = addBar(p, root, stats)
		rootStats[root] = &RootStat{Size: stats}
	}

	for p := range progress {
		root := strings.Split(p.Info.RelativePath, string(filepath.Separator))[0]

		totalProgress.IncrBy(int(p.Progress), time.Since(start))

		bars[root].IncrBy(int(p.Progress), time.Since(start))
		rootStats[root].Downloaded += p.Progress
	}
	p.Wait()
}
