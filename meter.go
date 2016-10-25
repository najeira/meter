package meter

import (
	"math"
	"sync"
	"time"
)

type Meter interface {
	Count() int64
	Rate() float64
	Mark(int64)
}

func New(duration time.Duration) Meter {
	return newRateMeter(duration, decayLinear)
}

func NewMeanLifetime(duration time.Duration) Meter {
	return newRateMeter(duration, decayMeanLifetime)
}

func NewHalfLife(duration time.Duration) Meter {
	return newRateMeter(duration, decayHalfLife)
}

func newRateMeter(duration time.Duration, strategy decayStrategy) Meter {
	return &rateMeter{
		count:     0,
		rate:      0,
		decayed:   time.Now().UnixNano(),
		duration:  duration.Nanoseconds(),
		storategy: strategy,
	}
}

type rateMeter struct {
	mu        sync.RWMutex
	count     int64
	rate      float64
	decayed   int64
	duration  int64
	storategy decayStrategy
}

var _ Meter = (*rateMeter)(nil)

func (m *rateMeter) Count() int64 {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.count
}

func (m *rateMeter) Rate() float64 {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.decay()
	return m.rate
}

func (m *rateMeter) decay() {
	now := time.Now().UnixNano()
	elapsed := now - m.decayed
	if elapsed > m.duration {
		m.rate = 0
	} else {
		fraction := (float64(elapsed) / float64(m.duration))
		ratio := m.storategy(fraction)
		m.rate = m.rate * ratio
	}
	m.decayed = now
}

func (m *rateMeter) Mark(count int64) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.count += count

	m.decay()
	m.rate += float64(count)
}

type decayStrategy func(float64) float64

func decayLinear(fraction float64) float64 {
	return 1 - fraction
}

func decayMeanLifetime(fraction float64) float64 {
	return math.Pow(math.E, -fraction)
}

func decayHalfLife(fraction float64) float64 {
	return math.Pow(0.5, fraction)
}
