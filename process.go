package main

import (
	"encoding/json"
	"fmt"
	"runtime"
	"stock/updater/depth"
	"stock/updater/depth/bs"
	"stock/updater/models"
)

func process(in chan []byte) (out chan []byte, done chan struct{}) {
	asksDepth := bs.NewDepth(depth.Options{Descending: false})
	bidsDepth := bs.NewDepth(depth.Options{Descending: true})

	_ = asksDepth
	_ = bidsDepth

	out = make(chan []byte, runtime.NumCPU())
	done = make(chan struct{}, 0)

	go func() {
		defer close(done)
		defer close(out)
		currentSnapShot := models.Data{}
		data := models.Data{}
		for {
			select {
			case info, ok := <-in:
				if !ok {
					return
				}
				_ = info

				err := json.Unmarshal(info, &data)
				if err != nil {
					fmt.Println(err)
					continue
				}

				switch data.Type {
				case models.Snapshot:
					asks := make([]models.Item, len(data.Payload.Asks))
					copy(asks, data.Payload.Asks)
					bids := make([]models.Item, len(data.Payload.Bids))
					copy(bids, data.Payload.Bids)
					asksDepth.Store(asks)
					bidsDepth.Store(bids)
					currentSnapShot = data
				case models.Update:
					for _, item := range data.Payload.Asks {
						asksDepth.Update(item)
					}
					for _, item := range data.Payload.Bids {
						bidsDepth.Update(item)
					}

					currentSnapShot.Payload.Asks = asksDepth.Load()
					currentSnapShot.Payload.Bids = bidsDepth.Load()

					body, err := json.Marshal(currentSnapShot)
					if err != nil {
						fmt.Println(err)
						continue
					}

					out <- body
				}
			}
		}
	}()

	return
}
