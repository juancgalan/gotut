package main

import (
	"fmt"
	"math/rand"
	"sync"
)

type Matrixf [][]float64
type Vectorf []float64

type FutureF64 struct {
	barrier    []chan struct{}
	data       float64
	resolved   bool
	mutex      *sync.RWMutex
	barrierMtx *sync.Mutex
}

func NewFutureF64() FutureF64 {
	return FutureF64{
		barrier:    make([]chan struct{}, 0),
		data:       0,
		resolved:   false,
		mutex:      &sync.RWMutex{},
		barrierMtx: &sync.Mutex{},
	}
}

func (f *FutureF64) set(val float64) {
	if f.resolved {
		return
	}
	f.mutex.Lock()
	f.data = val
	f.resolved = true
	for i := range f.barrier {
		close(f.barrier[i])
	}
	f.mutex.Unlock()
}

func (f *FutureF64) get() float64 {
	f.mutex.RLock()
	if f.resolved {
		return f.data
	}
	c := make(chan struct{})
	f.barrierMtx.Lock()
	f.barrier = append(f.barrier, c)
	f.barrierMtx.Unlock()
	f.mutex.RUnlock()
	<-c
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
	fx := make([]FutureF64, len(target))
	for i := range fx {
		fx[i] = NewFutureF64()
	}
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
					sum += system[id][j] * fx[j].get()
				}
				fx[id].set((target[id] - sum) / system[id][id])
			}
		}(c[i])
	}
	for i := range system {
		c[i%cpu] <- i
	}
	for i := range c {
		c[i] <- -1
	}
	for i := range fx {
		x[i] = fx[i].get()
	}
	return
}

func main() {
	vec := []float64{1, 2, 3, 4, 5}
	mtx := [][]float64{
		{1, 0, 0, 0, 0},
		{2, 7, 0, 0, 0},
		{3, 8, 13, 0, 0},
		{4, 9, 14, 19, 0},
		{5, 10, 15, 20, 25},
	}
	result := mmul(vec, mtx)
	e := parallelSolveTriang2(mtx, result, 2)
	for i := range e {
		fmt.Printf("%2.2f ", e[i])
		if e[i] != vec[i] {
			fmt.Printf("Incorrect cell value,\n got %v, \nexpected %v, \ntarget %v\n", e, vec, result)
		}
	}
}
