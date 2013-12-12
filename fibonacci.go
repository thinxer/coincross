package gocoins

import (
	"time"
)

type Fibonacci struct {
	a, b int64
}

func makeFibonacci() *Fibonacci {
	return &Fibonacci{
		a: 0,
		b: 1,
	}
}

func (f *Fibonacci) Next() int64 {
	f.a, f.b = f.b, f.a+f.b
	return f.a
}

func (f *Fibonacci) Prev() int64 {
	if f.a > 0 {
		f.a, f.b = f.b-f.a, f.a
	}
	return f.a
}

// If f returns true, repeat at give interval.
// Otherwise, sleep in seconds of Fibonacci series.
func FibonacciTimer(f func() bool, interval time.Duration) {
	fib := makeFibonacci()
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
