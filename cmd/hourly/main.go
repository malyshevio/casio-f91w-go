package main

import (
	"flag"
	"log"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/wav"
	"github.com/faiface/beep/speaker"

	//"github.com/MawKKe/casio-f91w-go/speaker"
	assets "github.com/MawKKe/casio-f91w-go/assets"
	"github.com/MawKKe/casio-f91w-go/ticker"
)

func main() {
	var debug bool
	flag.BoolVar(&debug, "debug", false, "beep on every new tenth second instead of every hour")
	flag.Parse()

	f, err := assets.Assets.Open("beepBeepNoTrailingDelay.wav")
	if err != nil {
		log.Fatal(err)
	}

	clipStreamer, format, err := wav.Decode(f)
	if err != nil {
		log.Fatal(err)
	}
	defer clipStreamer.Close()

	pred := ticker.HourChangePred

	if debug {
		pred = ticker.TenSecondChangePred
	}

	tickr := ticker.NewTicker(pred, 250*time.Millisecond)

	done := make(chan bool)

	go tickr.Run()

	log.Println("starting loop")
	for {
		select {
		case <-tickr.NotifyChannel():
			// A bit hacky but beep lib is ass. Doint Init() starts a goroutine
			// that polls continuously in the background even when nothing is
			// playing. The .Close() stops the goroutine, but requires re-Init()
			// before next play.
			speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
			speaker.Play(beep.Seq(clipStreamer, beep.Callback(func() {
				done <- true
			})))
			log.Printf("beep beep")
			<-done
			speaker.Close()
			clipStreamer.Seek(0)
		}
	}
}
