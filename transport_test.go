package main

import (
	"testing"
)

func TestCalculateSample1(t *testing.T) {
	f := Facts{Costs: [][]float64{
		{8, 3, 1},
		{4, 7, 4},
		{5, 2, 6},
	},
		Demands:  []float64{70, 60, 30},
		Supplies: []float64{30, 90, 50},
	}
	find(t, f)

}

func TestCalculateSample2(t *testing.T) {
	f := Facts{Costs: [][]float64{
		{2, 3, 4},
		{3, 2, 3},
		{2, 5, 2},
		{4, 1, 6},
	},
		Demands:  []float64{20, 20, 30, 10},
		Supplies: []float64{30, 40, 20},
	}
	find(t, f)
}

func TestCalculateSample3(t *testing.T) {
	f := Facts{Costs: [][]float64{
		{2, 3, 4},
		{3, 2, 3},
		{2, 5, 2},
		{4, 1, 6},
	},
		Demands:  []float64{20, 20, 30, 10},
		Supplies: []float64{30, 40, 20, 10},
	}
	if _, err := NewFactsWrapper(f); err != ErrorSLength {
		t.Fatalf("failed error responce: %v", err)
	}
}

func TestCalculateSample4(t *testing.T) {
	f := Facts{Costs: [][]float64{
		{8, 3, 1},
		{4, 7, 4},
		{5, 2, 6},
	},
		Demands:  []float64{70, 60, 30, 4, 5},
		Supplies: []float64{30, 90, 50},
	}
	if _, err := NewFactsWrapper(f); err != ErrorDLength {
		t.Fatalf("failed error responce: %v", err)
	}
}

func TestCalculateSample5(t *testing.T) {
	f := Facts{Costs: [][]float64{
		{8, 3, 1},
		{4, 7, 4},
		{5, 2, 6},
		{1, 5, 2},
		{1, 5, 2},
	},
		Demands:  []float64{70, 60, 30, 4, 5},
		Supplies: []float64{30, 90, 50},
	}
	find(t, f)
}

func TestCalculateSample6(t *testing.T) {
	f := Facts{Costs: [][]float64{
		{8, 3, 1, 6, 3},
		{4, 7, 4, 6, 3},
		{5, 2, 6, 3, 7},
		{1, 3, 3, 9, 2},
	},
		Demands:  []float64{70, 60, 30, 23},
		Supplies: []float64{30, 90, 94, 10, 104},
	}
	find(t, f)
}

func TestCalculateSample7(t *testing.T) {
	f := Facts{Costs: [][]float64{
		{2, 7, 8},
		{7, 2, 5},
		{5, 8, 1},
	},
		Demands:  []float64{40, 30, 40},
		Supplies: []float64{30, 30, 10},
	}
	find(t, f)
}
func TestCalculateSample8(t *testing.T) {
	f := Facts{Costs: [][]float64{
		{8, 3, 1, 6, 3},
		{4, 7, 4, 6, 3},
		{5, 2, 6, 3, 7},
		{1, 3, 3, 9, 2},
		{1, 3, 3, 9, 2},
		{1, 3, 3, 9, 2},
		{1, 3, 3, 9, 2},
		{1, 3, 3, 9, 2},
	},
		Demands:  []float64{70, 60, 30, 23, 70, 60, 30, 23},
		Supplies: []float64{30, 90, 94, 10, 54},
	}
	find(t, f)
}

func find(t *testing.T, f Facts) {
	fw, err := NewFactsWrapper(f)
	if err != nil {
		t.Fatalf("failed to create new FactsWrapper, error: %v", err)
	}
	if err := fw.Find(); err != nil {
		t.Fatalf("failed to find decision, error: %v", err)
	}
}
