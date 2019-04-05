package main

import (
	"math/rand"
	"sync"
)

type Matrixf [][]float64
type Vectorf []float64

type FutureFt64 struct {
	barrier  map[chan struct{}]struct{}
	data     float64
	resolved boolean
	mutex    *sync.Mutex
}

func NewFutureF64() FutureF64 {
	return FutureF64{
		barrier:  make(map[chan struct{}]struct{}),
		data:     0,
		resolved: false,
		mutex:    &sync.Mutex{},
	}
}

func (f *FutureFloat64) set(val float64) {
	if f.resolved {
		return
	}
	f.data = val
	f.resolved = true
	for i := range f.barrier {
		close(f.barrier[i])
	}
}

func (f *FutureFloat64) get() {
	if f.resolved {
		return f.data
	}
	c := make(chan struct{})
	f.barrier[c] = nil
	_ = <-c
	return f.data
}

func zeros(x, y int) (ans Matrixf) {
	ans = make([][]float64, x, x)
	for i := range ans {
		ans[i] = make([]float64, y, y)
		for j := range ans[i] {
			ans[i][j] = 0
		}
	}
	return
}

func zeroVec(x int) (ans Vectorf) {
	ans = make([]float64, x, x)
	for i := range ans {
		ans[i] = 0
	}
	return
}

func randMtx(x, y int, n float64) (ans Matrixf) {
	ans = make([][]float64, x, x)
	for i := range ans {
		ans[i] = make([]float64, y, y)
		for j := range ans[i] {
			ans[i][j] = n * rand.Float64()
		}
	}
	return
}

func randVec(x int, n float64) (ans Vectorf) {
	ans = make([]float64, x, x)
	for i := range ans {
		ans[i] = n * rand.Float64()
	}
	return
}

func mmul(v Vectorf, m Matrixf) (ans Vectorf) {
	ans = zeroVec(len(v))
	for i := range m {
		sum := 0.0
		for j := range v {
			sum += m[i][j] * v[j]
		}
		ans[i] = sum
	}
	return
}

func parallelMmul(v Vectorf, m Matrixf, cpu int) (ans Vectorf) {
	ans = zeroVec(len(v))
	var wg sync.WaitGroup
	slice := len(v) / cpu
	for i := 0; i < cpu; i += 1 {
		wg.Add(1)
		go func(from, to int) {
			defer wg.Done()
			for k := from; k < to; k += 1 {
				sum := 0.0
				for j := range v {
					sum += m[k][j] * v[j]
				}
				ans[k] = sum
			}
		}(i*slice, (i+1)*slice)
	}
	wg.Wait()
	return
}

func solveTriang(system Matrixf, target Vectorf) (x Vectorf) {
	x = zeroVec(len(target))
	for i := range system {
		sum := 0.0
		for j := 0; j < i; j += 1 {
			sum += system[i][j] * x[j]
		}
		x[i] = (target[i] - sum) / system[i][i]
	}
	return
}

func parallelSolveTriang(system Matrixf, target Vectorf, cpu int) (x Vectorf) {
	x = zeroVec(len(target))
	c := make([]chan int, cpu, cpu)
	for i := range c {
		c[i] = make(chan int)
		go func(idChan chan int) {
			for {
				id := <-idChan
				if id < 0 {
					return
				}
				sum := 0.0
				for j := 0; j < id; j += 1 {
					sum += system[id][j] * x[j]
				}
				x[id] = (target[id] - sum) / system[id][id]
			}
		}(c[i])
	}
	for i := range system {
		c[i%cpu] <- i
	}
	for i := range c {
		c[i] <- -1
	}
	return
}

func parallelSolveTriang2(system Matrixf, target Vectorf, cpu int) (x Vectorf) {
	x = zeroVec(len(target))
	c := make([]chan int, cpu, cpu)
	for i := range c {
		c[i] = make(chan int)
		go func(idChan chan int) {
			for {
				id := <-idChan
				if id < 0 {
					return
				}
				sum := 0.0
				for j := 0; j < id; j += 1 {
					sum += system[id][j] * x[j]
				}
				x[id] = (target[id] - sum) / system[id][id]
			}
		}(c[i])
	}
	for i := range system {
		c[i%cpu] <- i
	}
	for i := range c {
		c[i] <- -1
	}
	return
}

func main() {

}
