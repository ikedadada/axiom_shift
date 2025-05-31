package logic

import (
	"math/rand"
)

// RuleMatrix represents the matrix of rules that will be used in battles.
type RuleMatrix struct {
	matrix [][]float64
}

// NewRuleMatrix generates a new RuleMatrix based on a given seed.
func NewRuleMatrix(seed int64, size int) *RuleMatrix {
	rand.Seed(seed)
	matrix := make([][]float64, size)
	for i := range matrix {
		matrix[i] = make([]float64, size)
		for j := range matrix[i] {
			matrix[i][j] = rand.Float64()*2 - 1 // -1〜+1の範囲でランダム
		}
	}
	return &RuleMatrix{matrix: matrix}
}

// GetMatrix returns the underlying matrix.
func (r *RuleMatrix) GetMatrix() [][]float64 {
	return r.matrix
}
