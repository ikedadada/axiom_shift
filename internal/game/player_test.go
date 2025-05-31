package game

import "testing"

func TestNewCharacter(t *testing.T) {
	m := NewMatrix(2, 2)
	c := NewCharacter(m, 0.5)
	if c.MatrixState != m || c.GrowthRate != 0.5 {
		t.Error("NewCharacter did not set fields correctly")
	}
}

func TestCharacterUpdateMatrix(t *testing.T) {
	m := NewMatrix(2, 2)
	c := NewCharacter(m, 1.0)
	// 入力値0.0→インデックス0、要素[0][0]が大きく、他は微小成長
	c.UpdateMatrix(0.0)
	if m.data[0][0] != 1.0 {
		t.Errorf("Expected m.data[0][0]=1.0, got %v", m.data[0][0])
	}
	if m.data[0][1] != 0.1 || m.data[1][0] != 0.1 || m.data[1][1] != 0.1 {
		t.Errorf("Other elements should be 0.1, got %v", m.data)
	}
	// 入力値1.0→インデックス3、要素[1][1]が大きく、他は微小成長
	c.UpdateMatrix(1.0)
	if m.data[1][1] != 1.1 {
		t.Errorf("Expected m.data[1][1]=1.1, got %v", m.data[1][1])
	}
	if m.data[0][0] != 1.1 || m.data[0][1] != 0.2 || m.data[1][0] != 0.2 {
		t.Errorf("Other elements should be incremented by 0.1, got %v", m.data)
	}
}

func TestCharacterGetMatrix(t *testing.T) {
	m := NewMatrix(2, 2)
	c := NewCharacter(m, 1.0)
	if c.GetMatrix() != m {
		t.Error("GetMatrix did not return correct matrix")
	}
}
