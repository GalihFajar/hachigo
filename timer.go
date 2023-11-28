package main

import (
	"sync"
	"time"
)

type Timer struct {
	mu                sync.Mutex
	Time              byte
	TimerCallback     []func()
	TimerZeroCallback []func()
}

func (t *Timer) GetTime() byte {
	t.mu.Lock()
	time := t.Time
	defer t.mu.Unlock()
	return time
}

func (t *Timer) SetTime(in byte) {
	t.mu.Lock()
	t.Time = in
	defer t.mu.Unlock()
}

func (t *Timer) Decrement() {
	f := func() {
		t.mu.Lock()
		defer t.mu.Unlock()
		if t.Time > 0 {
			for _, c := range t.TimerCallback {
				c()
			}

			t.Time -= 1
		} else {
			for _, c := range t.TimerZeroCallback {
				c()
			}
		}
	}
	ticker := time.NewTicker(1 * (time.Second / 50))

	go func() {
		for range ticker.C {
			if destroy {
				ticker.Stop()
			}
			f()
		}
	}()
}
