package logic

import (
	"testing"
)

func TestNewRuleMatrix(t *testing.T) {
	rule := NewRuleMatrix(42, 2)
	if rule == nil || len(rule.matrix) != 2 || len(rule.matrix[0]) != 2 {
		t.Error("NewRuleMatrix did not create correct size matrix")
	}
}

func TestRuleMatrixGetMatrix(t *testing.T) {
	rule := NewRuleMatrix(42, 2)
	mat := rule.GetMatrix()
	if len(mat) != 2 || len(mat[0]) != 2 {
		t.Error("GetMatrix did not return correct size")
	}
}
