package test

import (
	"testing"

	"gonum.org/v1/gonum/mat"

	"github.com/jamestjsp/harold-go/harold"
)

func TestState(t *testing.T) {
	// Define matrices
	a := mat.NewDense(2, 2, []float64{0, 1, -4, -5})
	b := mat.NewDense(2, 1, []float64{0, 1})
	c := mat.NewDense(1, 2, []float64{1, 0})
	d := mat.NewDense(1, 1, []float64{1})

	// Create State instance
	state, err := harold.NewState(a, b, c, d, 0.1)
	if err != nil {
		t.Fatalf("Failed to create State: %v", err)
	}

	// Test SamplingPeriod
	if state.SamplingPeriod() != 0.1 {
		t.Errorf("Expected SamplingPeriod to be 0.1, got %v", state.SamplingPeriod())
	}

	// Test SamplingSet
	if state.SamplingSet() != "Z" {
		t.Errorf("Expected SamplingSet to be 'Z', got %v", state.SamplingSet())
	}

	// Test NumberOfStates
	if state.NumberOfStates() != 2 {
		t.Errorf("Expected NumberOfStates to be 2, got %v", state.NumberOfStates())
	}

	// Test NumberOfInputs
	if state.NumberOfInputs() != 1 {
		t.Errorf("Expected NumberOfInputs to be 1, got %v", state.NumberOfInputs())
	}

	// Test NumberOfOutputs
	if state.NumberOfOutputs() != 1 {
		t.Errorf("Expected NumberOfOutputs to be 1, got %v", state.NumberOfOutputs())
	}

	// Test Shape
	p, m := state.Shape()
	if p != 1 || m != 1 {
		t.Errorf("Expected Shape to be (1, 1), got (%v, %v)", p, m)
	}

	// Test Matrices
	a1, b1, c1, d1 := state.Matrices()
	if !mat.Equal(a, a1) || !mat.Equal(b, b1) || !mat.Equal(c, c1) || !mat.Equal(d, d1) {
		t.Errorf("Matrices do not match")
	}

	// Test ToArray
	if _, err := state.ToArray(); err == nil {
		t.Errorf("Expected error when calling ToArray on non-gain state")
	}

	// Test String
	expectedDesc := "State representation with sampling time: 0.100\n2 states, 1 inputs, and 1 outputs\n"
	if state.String() != expectedDesc {
		t.Errorf("Expected String to be %v, got %v", expectedDesc, state.String())
	}
}

func TestStateNegative(t *testing.T) {
	// Test with incompatible matrix dimensions between A and C
	a := mat.NewDense(2, 2, []float64{0, 1, -4, -5})
	b := mat.NewDense(2, 1, []float64{0, 1})
	c := mat.NewDense(1, 3, []float64{1, 0, 2}) // Incompatible with A (columns don't match)
	d := mat.NewDense(1, 1, []float64{1})

	_, err := harold.NewState(a, b, c, d, 0.1)
	if err == nil {
		t.Errorf("Expected error due to incompatible C matrix dimensions, but got nil")
	}

	// Test with incompatible matrix dimensions between A and B
	bBad := mat.NewDense(3, 1, []float64{0, 1, 2}) // Incompatible with A (rows don't match)
	_, err = harold.NewState(a, bBad, mat.NewDense(1, 2, []float64{1, 0}), d, 0.1)
	if err == nil {
		t.Errorf("Expected error due to incompatible B matrix dimensions, but got nil")
	}

	// Test static gain (A is nil) but no D matrix
	_, err = harold.NewState(nil, b, c, nil, 0.1)
	if err == nil {
		t.Errorf("Expected error due to missing D for static gain, but got nil")
	}

	// Test with nil B matrix (using compatible C matrix)
	goodC := mat.NewDense(1, 2, []float64{1, 0}) // Compatible with A (2 columns)
	_, err = harold.NewState(a, nil, goodC, d, 0.1)
	if err != nil {
		t.Errorf("Expected nil B to be valid, but got error: %v", err)
	}

	// Test with nil C matrix
	_, err = harold.NewState(a, b, nil, d, 0.1)
	if err != nil {
		t.Errorf("Expected nil C to be valid, but got error: %v", err)
	}
}
