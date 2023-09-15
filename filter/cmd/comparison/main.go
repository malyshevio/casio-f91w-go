package main

import (
	"log"
	"time"
)

func main() {
	t := time.NewTicker(time.Second)

	start := time.Now()
	stop := start.Add(45 * time.Second)

	log.Println("start")
loop:
	for {
		select {
		case tick := <-t.C:
			log.Printf("beep beep: %+v", tick)
			if tick.After(stop) {
				break loop
			}
		}
	}
	log.Println("end")
}
