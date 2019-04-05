package main

import "math"

func min(data [][]float64) (minI, minJ int, val float64) {
	val = math.MaxFloat64
	for i, item := range data {
		for j, cell := range item {
			if cell < val {
				val = cell
				minI = i
				minJ = j
			}
		}
	}
	return
}
func max(data [][]float64) (minI, minJ int, val float64) {
	val = -math.MaxFloat64
	for i, item := range data {
		for j, cell := range item {
			if cell > val && cell != empty {
				val = cell
				minI = i
				minJ = j
			}
		}
	}
	return
}

func sum(data []float64) (sum float64) {
	for _, num := range data {
		sum += num
	}
	return
}

func makeTable(size1, size2 int) (resp [][]float64) {
	resp = make([][]float64, size1)
	for i := range resp {
		resp[i] = make([]float64, size2)
		reset(resp[i])
	}
	return
}

func makeSlice(size int) (resp []float64) {
	resp = make([]float64, size)
	reset(resp)
	return
}

func reset(a []float64) {
	for i := range a {
		a[i] = empty
	}
}
