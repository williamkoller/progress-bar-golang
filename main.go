package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/vbauerster/mpb/v8"
	"github.com/vbauerster/mpb/v8/decor"
)

func main() {
	var wg sync.WaitGroup
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	p := mpb.New(mpb.WithWaitGroup(&wg))

	taskCount := 3
	bars := make([]*mpb.Bar, taskCount)
	totals := make([]int, taskCount)

	for i := 0; i < taskCount; i++ {
		total := rnd.Intn(100) + 50
		totals[i] = total

		bars[i] = p.AddBar(int64(total),
			mpb.PrependDecorators(
				decor.Name(fmt.Sprintf("Task #%d:", i+1)),
				decor.CountersNoUnit("%d/%d"),
			),
			mpb.AppendDecorators(
				decor.Percentage(decor.WC{W: 5}),
				decor.OnComplete(decor.Name(" ✔️ "), " ✅"),
			),
		)
	}

	for i := 0; i < taskCount; i++ {
		wg.Add(1)
		go func(id, total int, bar *mpb.Bar) {
			defer wg.Done()

			for j := 0; j < total; j++ {
				bar.Increment()
				time.Sleep(time.Duration(rnd.Intn(30)+10) * time.Millisecond)
			}
		}(i, totals[i], bars[i])
	}

	p.Wait()
}
