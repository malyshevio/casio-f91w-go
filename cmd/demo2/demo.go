package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/MawKKe/casio-f91w-go/ticker"

	"github.com/faiface/beep"
	"github.com/faiface/beep/generators"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/wav"
)

func usage() {
	fmt.Printf("usage: %s freq\n", os.Args[0])
	fmt.Println("where freq must be a float between 1 and 24000")
	fmt.Println("24000 because samplerate of 48000 is hardcoded")
}

const (
	BURST_INTERVAL      = 600
	BURST_HALF_LENGTH   = 150
	BURST_HALF_INTERVAL = 100
)

// thereabouts..
var casioFreq int = 8000.0

func CasioF91HourlyBeep(sr beep.SampleRate) []beep.Streamer {
	burstHalf := sr.N(BURST_HALF_LENGTH * time.Millisecond)

	// test different generators. It is likely sawtooth or rectangular but who knows?
	sound, err := generators.SinTone(sr, casioFreq)
	if err != nil {
		panic(err)
	}

	sounds := []beep.Streamer{
		beep.Take(burstHalf, sound),
	}
	return sounds
}

type BeepState int

const (
	FirstBeep BeepState = iota
	InterPause
	SecondBeep
	PostPause
	BeepStateEnd
)

type BeepStrat int

const (
	Default BeepStrat = iota
	OnlySingleBeep
	OnlyDoubleBeep
)

func BeepWithContext(ctx context.Context) {
	log.Println("enter")
	f, err := os.Open("assets/beep.wav")
	if err != nil {
		log.Fatal(err)
	}

	sound, format, err := wav.Decode(f)
	if err != nil {
		log.Fatal(err)
	}
	sr := beep.SampleRate(format.SampleRate)

	//iter GetBeeper(sr, Default)

	burstInterval := sr.N(BURST_INTERVAL * time.Millisecond)
	burstHalfInterval := sr.N(BURST_HALF_INTERVAL * time.Millisecond)

	log.Println("len wav:", sound.Len())
	log.Println("len half interval:", burstHalfInterval)
	log.Println("len interval:", burstInterval)
	log.Println("tot", sound.Len()*2+burstHalfInterval+burstInterval)

	state := FirstBeep

	strat := Default

	iter := func() beep.Streamer {
		switch state {
		case FirstBeep:
			log.Println(time.Now())
			sound.Seek(0)
			if strat == OnlySingleBeep {
				state = BeepStateEnd
			} else {
				state = InterPause
			}
			return sound
		case InterPause:
			state = SecondBeep
			return beep.Silence(burstHalfInterval)
		case SecondBeep:
			sound.Seek(0)
			if strat == OnlyDoubleBeep {
				state = BeepStateEnd
			} else {
				state = PostPause
			}
			return sound
		case PostPause:
			state = BeepStateEnd
			return beep.Silence(burstInterval)
		default:
			select {
			case <-ctx.Done():
				return nil
			default:
				if strat == Default {
					state = FirstBeep
				}
				return beep.Silence(0)
			}
		}
	}

	sounds := beep.Iterate(iter)

	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	defer speaker.Close() // gotta stop it, unless it stays in the background consuming ~5%

	closed := make(chan struct{})

	log.Println("playing")
	speaker.Play(beep.Seq(sounds, beep.Callback(func() { closed <- struct{}{} })))

	log.Println("waiting")
	select {
	case <-closed:
		break
	case <-ctx.Done():
		break
	}
	speaker.Close()

	log.Println("finished, closed")
}

func main() {
	var debug bool
	flag.BoolVar(&debug, "debug", false, "beep on every new tenth second instead of every hour")
	flag.Parse()

	pred := ticker.HourChangePred
	if debug {
		pred = ticker.TenSecondChangePred
	}

	tickr := ticker.NewTicker(pred, 250*time.Millisecond)

	ctx, cancel := context.WithTimeout(context.Background(), time.Hour)
	defer cancel()

	log.Println("launch")

	// casio.BeepWithContext(SingleBeep, ctx)
	// casio.BeepWithContext(DoubleBeep, ctx)

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		tickr.RunWithContext(ctx)
	}()

loop:
	for {
		select {
		case <-tickr.NotifyChannel():
			wg.Add(1)
			go func() {
				defer wg.Done()
				BeepWithContext(ctx)
			}()
		case <-ctx.Done():
			break loop
		}
	}
	wg.Wait()
	log.Println("done")

}
