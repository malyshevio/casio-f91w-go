package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"time"

	//"github.com/MawKKe/casio-f91w-go/speaker"
	assets "github.com/MawKKe/casio-f91w-go/assets"
	filter "github.com/MawKKe/casio-f91w-go/filter"

	"github.com/ebitengine/oto/v3"
)

func clamp(lo, value, hi float64) float64 {
	if value < lo {
		return 0.0
	}
	if value > hi {
		return 1.0
	}
	return value
}

func main() {
	var debug bool
	var volume float64
	flag.BoolVar(&debug, "debug", false, "debug mode - beep every second")
	flag.Float64Var(&volume, "volume", 1.0, "set volume (range 0.0 ... 1.0)")
	flag.Parse()

	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	log.Println("Starting..")

	// The oto documentation says the Player.SetVolume() expects the volume in range 0..1
	// yet it accepts any value, possibly distorting the audio. Go figure.
	volume = clamp(0.0, volume, 1.0)

	// we use raw S16LE PCM so we don't have to worry about decoding
	pcmBytes, err := assets.Assets.ReadFile("beepBeepNoTrailingDelay.pcm")
	if err != nil {
		log.Fatal(err)
	}

	op := &oto.NewContextOptions{
		SampleRate:   44100,
		ChannelCount: 1,
		Format:       oto.FormatSignedInt16LE,
	}

	// NOTE: max one context per program
	otoCtx, readyChan, err := oto.NewContext(op)
	if err != nil {
		panic("oto.NewContext failed: " + err.Error())
	}

	// It might take a bit for the hardware audio devices to be ready, so we wait on the channel.
	<-readyChan

	log.Println("Audio context initialized")

	// NOTE: player is Pause()'d by default
	player := otoCtx.NewPlayer(bytes.NewReader(pcmBytes))

	defer func() {
		if err := player.Close(); err != nil {
			panic(fmt.Errorf("player.Close failed: %q", err))
		}
	}()

	player.SetVolume(volume)

	var floorFunc filter.RoundingFunc
	var tickerDuration time.Duration

	if debug {
		log.Println("Enabling beep on every ten seconds (debug)")
		floorFunc = filter.FloorToTensOfSeconds
		tickerDuration = 100 * time.Millisecond
	} else {
		log.Println("Enabling beep on every hour")
		floorFunc = filter.FloorToHour
		tickerDuration = time.Second
	}

	ticker := time.NewTicker(tickerDuration)
	tf := filter.NewTickerFilter(floorFunc)

	log.Println("Entering loop")

	if err := otoCtx.Suspend(); err != nil {
		panic(err)
	}

loop:
	for {
		var ts time.Time

		select {
		case ts = <-ticker.C:
			if !tf.HasChanged(ts) {
				continue loop
			}
		}
		if err := otoCtx.Resume(); err != nil {
			panic(err)
		}

		log.Printf("beep beep: +%v", ts.UTC())

		// Play starts playing the sound and returns without waiting for it (Play() is async).
		player.Play()

		// We can wait for the sound to finish playing using something like this
		for player.IsPlaying() {
			time.Sleep(time.Millisecond)
		}

		if _, err := player.Seek(0, io.SeekStart); err != nil {
			panic(fmt.Errorf("player.Seek() error: %q", err))
		}

		// suspend to avoid wasting 5% cpu while doing nothing...
		if err := otoCtx.Suspend(); err != nil {
			panic(err)
		}
	}
}
