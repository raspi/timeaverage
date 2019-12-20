package timeaverage

import (
	"fmt"
	"sync"
	"time"
)

type measurement struct {
	T time.Time // When measurement was taken
	V float64   // Value of taken measurement
}

// function which returns measurement
type measurementFunc func() (float64, error)

type TimeAverage struct {
	window          time.Duration // for how long measurements are kept (for example 5 minutes)
	rate            time.Duration // how often measurement is taken (for example every 500 ms)
	measurements    []measurement
	measurementFunc measurementFunc // function doing the measurement
	ticker          *time.Ticker    // ticker triggering measurement function
	running         bool            // is this running?
	average         float64         // calculated average value
	SampleCount     uint            // how many measurements in .measurements
	mu              sync.RWMutex    // mutex for locking while adding measurements and calculating averages
	Newest          time.Time       // what's the newest entry in .measurements
	Oldest          time.Time       // what's the oldest entry in .measurements (usually .measurements[0].T)
}

func New(timeWindow, rate time.Duration, initialavg float64, samplerFunc measurementFunc) *TimeAverage {
	return &TimeAverage{
		window:          timeWindow,
		rate:            rate,
		running:         false,
		measurementFunc: samplerFunc,
		average:         initialavg,
		Oldest:          time.Now(),
		Newest:          time.Now(),
		SampleCount:     0,
	}
}

func (s *TimeAverage) Stop() {
	s.ticker.Stop()
}

func (s *TimeAverage) Start() {
	if s.running {
		panic(fmt.Errorf(`err: Stop() was not called`))
	}

	s.measurements = []measurement{}
	s.SampleCount = 0

	s.ticker = time.NewTicker(s.rate)

	// Add first measurement
	s.addSample()

	go func() {
		for {
			select {
			case <-s.ticker.C:
				// Add a new measurement at each tick
				s.addSample()
			default:

			}
		}
	}()

	s.running = true
}

func (s *TimeAverage) addSample() {
	// Call sampler func for taking measurement
	v, err := s.measurementFunc()
	if err != nil {
		panic(err)
	}

	s.mu.Lock()
	now := time.Now()
	s.Newest = now
	s.measurements = append(s.measurements, measurement{T: now, V: v})
	s.updateAverage()
	s.SampleCount = uint(len(s.measurements))
	s.mu.Unlock()

}

func (s *TimeAverage) updateAverage() {
	var newSamples []measurement

	// Drop old measurement(s)
	for _, v := range s.measurements {
		if v.T.Before(time.Now().Add(-s.window)) {
			continue
		}

		newSamples = append(newSamples, v)
	}

	s.measurements = newSamples

	// Calculate average
	l := len(newSamples)

	tot := 0.0

	for idx, v := range newSamples {
		if idx == 0 {
			s.Oldest = v.T
		}

		tot += v.V
	}

	s.average = tot / float64(l)
}

func (s *TimeAverage) Average() float64 {
	if !s.running {
		panic(fmt.Errorf(`err: Start() was not called`))
	}

	return s.average
}
