package gocoins

import (
	"time"
)

type fibonacci struct {
	a, b int64
}

func newFibonacci() *fibonacci {
	return &fibonacci{
		a: 0,
		b: 1,
	}
}

func (f *fibonacci) Next() int64 {
	f.a, f.b = f.b, f.a+f.b
	return f.a
}

func (f *fibonacci) Prev() int64 {
	if f.a > 0 {
		f.a, f.b = f.b-f.a, f.a
	}
	return f.a
}

// If f returns true, repeat at give interval.
// Otherwise, sleep in seconds of fibonacci series.
func fibonacciTimer(f func() bool, interval time.Duration) {
	fib := newFibonacci()
	for {
		start := time.Now()
		r := f()
		dur := time.Now().Sub(start)
		if r {
			time.Sleep(interval - dur)
			fib.Prev()
		} else {
			time.Sleep(time.Duration(fib.Next()) * time.Second)
		}
	}
}
