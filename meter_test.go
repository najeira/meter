package meter

import (
	"testing"
	"time"
)

func TestMeterLinear(t *testing.T) {
	runTestMeter(t, New(time.Millisecond*1))
}

func TestMeterMeanLifetime(t *testing.T) {
	runTestMeter(t, NewMeanLifetime(time.Millisecond*1))
}

func TestMeterHalfLife(t *testing.T) {
	runTestMeter(t, NewHalfLife(time.Millisecond*1))
}

func runTestMeter(t *testing.T, meter Meter) {
	var count int64 = 0
	for i := 0; i < 10; i++ {
		meter.Mark(10)
		count += 10

		if c := meter.Count(); c != count {
			t.Errorf("Count %d, expect %d", c, count)
		}
		if r := meter.Rate(); r > float64(count) {
			t.Errorf("Rate %f, expect %d", r, count)
		}
	}
}

func TestMeterLinearAfterDuration(t *testing.T) {
	m := New(time.Millisecond*1).(*rateMeter)
	m.Mark(10)
	m.decayed = time.Now().Add(time.Millisecond*-1).UnixNano()
	if r := m.Rate(); r > 0 {
		t.Errorf("Rate %f, expect %d", r, 0)
	}
}
