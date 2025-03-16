package harold

import (
	"errors"
	"fmt"

	"gonum.org/v1/gonum/mat"
)

type State struct {
	A, B, C, D *mat.Dense
	dt         float64
	isSISO     bool
	isGain     bool
	isStable   bool
	poles      []complex128 // Add poles field
	zeros      []complex128 // Add zeros field
}

func NewState(a, b, c, d *mat.Dense, dt float64) (*State, error) {
	// Validate arguments
	var br, cc int
	var isSISO bool

	if b != nil {
		br, _ = b.Dims()
	}
	if c != nil {
		_, cc = c.Dims()
	}

	if a != nil {
		ar, ac := a.Dims()
		if b != nil && (ar != br) {
			return nil, errors.New("'A' and 'B' row dimensions do not match")
		}
		if c != nil && (ac != cc) {
			return nil, errors.New("'A' and 'C' column dimensions do not match")
		}
		isSISO = b != nil && c != nil && br == 1 && cc == 1
	} else {
		// Static gain check
		if d == nil {
			return nil, errors.New("'D' matrix must be provided for static gain models")
		}
		dr, dc := d.Dims()
		isSISO = dr == 1 && dc == 1
	}

	s := &State{
		A:      a,
		B:      b,
		C:      c,
		D:      d,
		dt:     dt,
		isSISO: isSISO,
		isGain: a == nil,
		poles:  []complex128{}, // Initialize empty slice
		zeros:  []complex128{}, // Initialize empty slice
	}

	s.recalc()
	return s, nil
}

func (s *State) recalc() {
	if s.isGain {
		s.isStable = true
		s.poles = []complex128{} // Clear poles for gain models
		s.zeros = []complex128{} // Clear zeros for gain models
	} else {
		// Calculate poles and zeros and stability
		s.poles = calculatePoles(s.A)
		s.zeros = calculateZeros(s)
		s.isStable = s.checkStability()
	}
}

// Helper functions for calculating poles
func calculatePoles(a *mat.Dense) []complex128 {
	// Estimate the eigenvalues of A to find poles
	var eigen mat.Eigen
	ok := eigen.Factorize(a, mat.EigenRight)
	if !ok {
		// Factorization failed
		return nil
	}

	// Get the dimensions of the original matrix
	r, _ := a.Dims()

	// Create a slice to hold the eigenvalues
	eigenvalues := make([]complex128, r)

	// Compute and store the eigenvalues
	eigen.Values(eigenvalues)

	return eigenvalues
}

// Helper functions for calculating zeros
func calculateZeros(s *State) []complex128 {
	// This would implement the transmission_zeros function from the Python code
	// For now, return an empty slice as a placeholder
	return []complex128{}
}

func (s *State) checkStability() bool {
	// Check stability based on poles
	// ...existing code...
	return true
}

func (s *State) SamplingPeriod() float64 {
	return s.dt
}

func (s *State) SamplingSet() string {
	if s.dt == 0 {
		return "R"
	}
	return "Z"
}

func (s *State) NumberOfStates() int {
	if s.A == nil {
		return 0
	}
	r, _ := s.A.Dims()
	return r
}

func (s *State) NumberOfInputs() int {
	if s.B == nil {
		return 0
	}
	_, c := s.B.Dims()
	return c
}

func (s *State) NumberOfOutputs() int {
	if s.C == nil {
		return 0
	}
	r, _ := s.C.Dims()
	return r
}

func (s *State) Shape() (int, int) {
	return s.NumberOfOutputs(), s.NumberOfInputs()
}

func (s *State) Matrices() (*mat.Dense, *mat.Dense, *mat.Dense, *mat.Dense) {
	return s.A, s.B, s.C, s.D
}

func (s *State) ToArray() (*mat.Dense, error) {
	if !s.isGain {
		return nil, errors.New("only static gain models can be converted to arrays")
	}
	return s.D, nil
}

func (s *State) String() string {
	desc := fmt.Sprintf("State representation with sampling time: %.3f\n", s.dt)
	if s.isGain {
		desc += fmt.Sprintf("%dx%d Static Gain\n", s.NumberOfOutputs(), s.NumberOfInputs())
	} else {
		desc += fmt.Sprintf("%d states, %d inputs, and %d outputs\n", s.NumberOfStates(), s.NumberOfInputs(), s.NumberOfOutputs())
	}
	return desc
}
