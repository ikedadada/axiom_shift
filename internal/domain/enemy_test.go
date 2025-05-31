package domain

import "testing"

func TestNewEnemy(t *testing.T) {
	m := NewMatrix(2, 2)
	e := NewEnemy("test", m, 0.5)
	if e.Name != "test" || e.Matrix != m || e.GrowthRate != 0.5 {
		t.Error("NewEnemy did not set fields correctly")
	}
}

func TestEnemyGrow(t *testing.T) {
	m := NewMatrix(2, 2)
	e := NewEnemy("test", m, 1.0)
	rule := NewMatrix(2, 2)
	rule.Data[0][0] = 1.0
	rule.Data[0][1] = -1.0
	rule.Data[1][0] = 0.0
	rule.Data[1][1] = 2.0
	e.Grow(0.0, rule)
	if m.Data[0][0] != 1.0 {
		t.Errorf("Grow did not update input element: got %v", m.Data[0][0])
	}
	if m.Data[1][1] != 1.1 {
		t.Errorf("Grow did not update max element: got %v", m.Data[1][1])
	}
	if m.Data[0][1] != 0.1 || m.Data[1][0] != 0.1 {
		t.Errorf("Other elements should be 0.1, got %v", m.Data)
	}
	e.Grow(1.0, rule)
	if m.Data[1][1] != 3.1 {
		t.Errorf("Grow did not update max element twice: got %v", m.Data[1][1])
	}
}

func TestEnemyGetMatrix(t *testing.T) {
	m := NewMatrix(2, 2)
	e := NewEnemy("test", m, 1.0)
	if e.GetMatrix() != m {
		t.Error("GetMatrix did not return correct matrix")
	}
}
