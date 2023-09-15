package main

import (
	"log"
	"time"

	filter "github.com/MawKKe/casio-f91w-go/filter"
)

func main() {
	tf := filter.NewTickerFilter(filter.FloorToTensOfSeconds)
	t := time.NewTicker(10 * time.Millisecond)

	start := time.Now()
	stop := start.Add(45 * time.Second)

	log.Println("start")

	// not strictly necessary but demonstrates how to feed the initial "previous" time.Time
	// tf.Init(time.Now())
loop:
	for {
		select {
		case tick := <-t.C:
			if tf.HasChanged(tick) {
				log.Printf("beep beep: %+v", tick)
			}
			if tick.After(stop) {
				break loop
			}
		}
	}
	log.Println("end")
}
