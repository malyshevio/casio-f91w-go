package ticker_filter_test

import (
	"reflect"
	"testing"
	"time"

	filter "github.com/MawKKe/casio-f91w-go/filter"
)

func TestFloorToSeconds(t *testing.T) {
	in := time.Date(2023, 9, 10, 13, 36, 14, 12345, time.UTC)
	expect := time.Date(2023, 9, 10, 13, 36, 14, 0, time.UTC)

	if res := filter.FloorToSeconds(in); res != expect {
		t.Fatalf("Expected: %q, got: %q", expect, res)
	}
}

func TestFloorToTenSeconds(t *testing.T) {
	in := time.Date(2023, 9, 10, 13, 36, 14, 12345, time.UTC)
	expect := time.Date(2023, 9, 10, 13, 36, 10, 0, time.UTC)

	if res := filter.FloorToTensOfSeconds(in); res != expect {
		t.Fatalf("Expected: %q, got: %q", expect, res)
	}
}

func TestFloorToHour(t *testing.T) {
	in := time.Date(2023, 9, 10, 13, 36, 14, 12345, time.UTC)
	expect := time.Date(2023, 9, 10, 13, 0, 0, 0, time.UTC)

	if res := filter.FloorToHour(in); res != expect {
		t.Fatalf("Expected: %q, got: %q", expect, res)
	}
}

func TestWithMockTicker(t *testing.T) {
	start := time.Date(2023, 9, 10, 13, 36, 14, 123456789, time.UTC)

	expect := []time.Time{
		time.Date(2023, 9, 10, 13, 36, 14, 0, time.UTC),
		time.Date(2023, 9, 10, 13, 36, 14, 0, time.UTC),
		time.Date(2023, 9, 10, 13, 36, 14, 0, time.UTC),
		time.Date(2023, 9, 10, 13, 36, 15, 0, time.UTC),
	}

	for i := 0; i < len(expect); i++ {
		dur := time.Duration(i*350) * time.Millisecond
		//result = append(result, start.Add(dur))
		if tmp := filter.FloorToSeconds(start.Add(dur)); tmp != expect[i] {
			t.Fatalf("Expected %q, got: %q", expect[i], tmp)
		}
	}
}

func TestTickerFilter(t *testing.T) {
	tf := filter.NewTickerFilter(filter.FloorToSeconds)
	start := time.Date(2023, 9, 10, 13, 36, 14, 123456789, time.UTC)

	var res []bool
	res = append(res, tf.HasChanged(start)) // first one, always false
	res = append(res, tf.HasChanged(start.Add(350*time.Millisecond)))
	res = append(res, tf.HasChanged(start.Add(700*time.Millisecond)))
	res = append(res, tf.HasChanged(start.Add(1050*time.Millisecond)))

	expect := []bool{false, false, false, true}

	if !reflect.DeepEqual(res, expect) {
		t.Fatalf("expect: %v, got: %v", expect, res)
	}
}

func TestTickerFilterInit(t *testing.T) {
	tf := filter.NewTickerFilter(filter.FloorToSeconds)
	start := time.Date(2023, 9, 10, 13, 36, 14, 123456789, time.UTC)

	if res := tf.HasChanged(start.Add(time.Second)); res {
		t.Fatal("Expected first call to HasChanged() to be false")
	}

	tf.Init(start)

	// this value is still in the same "seconds", so no change detected even after Init().
	if res := tf.HasChanged(start.Add(100 * time.Millisecond)); res {
		t.Fatal("Expected first call to HasChanged() (after Init()) to be false")
	}

	tf.Init(start)

	// the prev value is no longer zero and is in fact in the previous "seconds"
	// with respect to the time now queried.
	if res := tf.HasChanged(start.Add(time.Second)); !res {
		t.Fatal("Expected first call to HasChanged() (after Init()) to be true")
	}

}

func TestTickerFilterReset(t *testing.T) {
	tf := filter.NewTickerFilter(filter.FloorToSeconds)
	first := time.Date(2023, 9, 10, 13, 36, 14, 123456789, time.UTC)
	second := first.Add(time.Second)

	if res := tf.HasChanged(first); res {
		t.Fatal("Expected first call to HasChanged(first) to be false (step 1/3)")
	}

	if res := tf.HasChanged(second); !res {
		t.Fatal("Expected second call to HasChanged(second) to be true")
	}

	tf.Reset()
	if res := tf.HasChanged(first); res {
		t.Fatal("Expected third call to HasChanged(first) (after Reset()) to be false")
	}
}
