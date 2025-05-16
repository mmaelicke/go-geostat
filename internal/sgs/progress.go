package sgs

import (
	"fmt"
	"sync"
	"time"

	"github.com/vbauerster/mpb/v8"
	"github.com/vbauerster/mpb/v8/decor"
)

type Progress struct {
	id         int
	current    int
	total      int
	fitTime    time.Duration
	interpTime time.Duration
}

type progressTracker struct {
	progress    *mpb.Progress
	bars        map[int]*mpb.Bar
	barMux      sync.Mutex
	processChan chan Progress
}

func newProgressTracker() *progressTracker {
	pt := &progressTracker{
		progress:    mpb.New(),
		bars:        make(map[int]*mpb.Bar),
		processChan: make(chan Progress, 20),
	}

	go pt.manage()
	return pt
}

func (pt *progressTracker) manage() {
	for prog := range pt.processChan {
		pt.barMux.Lock()
		if _, exists := pt.bars[prog.id]; !exists {
			pt.bars[prog.id] = pt.progress.AddBar(int64(prog.total),
				mpb.PrependDecorators(
					decor.Name(fmt.Sprintf("Sim %3d:", prog.id)),
					decor.Percentage(decor.WCSyncSpace),
				),
				mpb.AppendDecorators(
					decor.EwmaETA(decor.ET_STYLE_GO, 30),
					decor.AverageSpeed(decor.CountersNoUnit, "%.2f"),
				),
			)
		}
		pt.bars[prog.id].Increment()
		pt.barMux.Unlock()
	}
}

func (pt *progressTracker) send(p Progress) {
	if pt != nil {
		pt.processChan <- p
	}
}

func (pt *progressTracker) close() {
	if pt != nil {
		close(pt.processChan)
		pt.progress.Wait()
	}
}
