package gocoins

type Fibonacci struct {
	a, b int64
}

func MakeFibonacci() *Fibonacci {
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
