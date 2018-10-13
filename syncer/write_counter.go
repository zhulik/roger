package syncer

import "io"

type WriteCounter struct {
	Writer   io.Writer
	Progress chan<- int64
	Total    int64
}

func (wc *WriteCounter) Write(p []byte) (int, error) {
	n, err := wc.Writer.Write(p)
	wc.Total += int64(n)
	wc.Progress <- int64(n)
	return n, err
}
