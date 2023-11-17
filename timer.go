package main

import (
	"fmt"
	"time"
)

type Timer struct {
	Time          byte
	TimerCallback []func()
}

func (t *Timer) Decrement() {
	f := func() {
		if t.Time >= 60 {

			for _, c := range t.TimerCallback {
				c()
			}

			t.Time -= 60
		} else {
			t.Time = 0

		}
	}
	ticker := time.NewTicker(1 * time.Second)

	go func() {
		for range ticker.C {
			f()
			fmt.Println("tick:", t.Time)
		}
	}()
}
