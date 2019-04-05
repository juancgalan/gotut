package main

import "testing"

func TestZeros(t *testing.T) {

	e := zeros(3, 3)
	if len(e) != 3 {
		t.Errorf("Incorrect columns size, got %d, expected %d", len(e), 3)
	}
	if len(e[0]) != 3 {
		t.Errorf("Incorrect rows size, got %d, expected %d", len(e[0]), 3)
	}
	for i := range e {
		for j := range e[i] {
			if e[i][j] != 0 {
				t.Errorf("Incorrect cell value, got %f, expected %d", e[i][j], 0)
			}
		}
	}
}

func TestZeroVec(t *testing.T) {

	e := zeroVec(3)
	if len(e) != 3 {
		t.Errorf("Incorrect vector columns size, got %d, expected %d", len(e), 3)
	}
	for i := range e {
		if e[i] != 0 {
			t.Errorf("Incorrect cell value, got %f, expected %d", e[i], 0)
		}
	}

}

func TestMmul(t *testing.T) {

	vec := []float64{1, 2, 3, 4, 5}
	mtx := [][]float64{
		{1, 6, 11, 16, 21},
		{2, 7, 12, 17, 22},
		{3, 8, 13, 18, 23},
		{4, 9, 14, 19, 24},
		{5, 10, 15, 20, 25},
	}
	exp := []float64{215, 230, 245, 260, 275}
	e := mmul(vec, mtx)
	for i := range e {
		if e[i] != exp[i] {
			t.Errorf("Incorrect cell value, got %v, expected %v", e, exp)
		}
	}
}

func TestParallelMmul(t *testing.T) {

	vec := []float64{1, 2, 3, 4, 5}
	mtx := [][]float64{
		{1, 6, 11, 16, 21},
		{2, 7, 12, 17, 22},
		{3, 8, 13, 18, 23},
		{4, 9, 14, 19, 24},
		{5, 10, 15, 20, 25},
	}
	exp := []float64{215, 230, 245, 260, 275}
	e := parallelMmul(vec, mtx, 1)
	for i := range e {
		if e[i] != exp[i] {
			t.Errorf("Incorrect cell value, got %v, expected %v", e, exp)
		}
	}
}

func TestSolveTriang(t *testing.T) {

	vec := []float64{1, 2, 3, 4, 5}
	mtx := [][]float64{
		{1, 0, 0, 0, 0},
		{2, 7, 0, 0, 0},
		{3, 8, 13, 0, 0},
		{4, 9, 14, 19, 0},
		{5, 10, 15, 20, 25},
	}
	result := mmul(vec, mtx)
	e := solveTriang(mtx, result)
	for i := range e {
		if e[i] != vec[i] {
			t.Errorf("Incorrect cell value, got %v, expected %v, target %v", e, vec, result)
		}
	}
}

func TestParallelSolveTriang(t *testing.T) {

	vec := []float64{1, 2, 3, 4, 5}
	mtx := [][]float64{
		{1, 0, 0, 0, 0},
		{2, 7, 0, 0, 0},
		{3, 8, 13, 0, 0},
		{4, 9, 14, 19, 0},
		{5, 10, 15, 20, 25},
	}
	result := mmul(vec, mtx)
	e := parallelSolveTriang(mtx, result, 2)
	for i := range e {
		if e[i] != vec[i] {
			t.Errorf("Incorrect cell value, got %v, expected %v, target %v", e, vec, result)
		}
	}
}

func TestParallelSolveTriang2(t *testing.T) {

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
		if e[i] != vec[i] {
			t.Errorf("Incorrect cell value, got %v, expected %v, target %v", e, vec, result)
		}
	}
}

func BenchmarkMmul(b *testing.B) {
	size := 3000 // 30000 for 1GB matrix
	vec := randVec(size, 1000.0)
	mtx := randMtx(size, size, 1000.0)
	b.ResetTimer()
	_ = mmul(vec, mtx)
}

func BenchmarkParallelMmul(b *testing.B) {
	size := 3000 // 30000 for 1GB matrix
	vec := randVec(size, 1000.0)
	mtx := randMtx(size, size, 1000.0)
	b.ResetTimer()
	_ = parallelMmul(vec, mtx, 2)
}
