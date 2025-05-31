package domain

import (
	"math/rand"
)

// RuleMatrix represents the matrix of rules that will be used in battles.
type RuleMatrix struct {
	*Matrix
}

// NewRuleMatrix generates a new RuleMatrix based on a given seed.
func NewRuleMatrix(seed int64, size int) *RuleMatrix {
	r := rand.New(rand.NewSource(seed))
	data := make([][]float64, size)
	for i := range data {
		data[i] = make([]float64, size)
		for j := range data[i] {
			data[i][j] = r.Float64()*2 - 1 // -1〜+1の範囲でランダム
		}
	}
	matrix := NewMatrix(data)
	return &RuleMatrix{
		Matrix: matrix,
	}
}
