package domain

import "testing"

func TestNewPlayer(t *testing.T) {
	m := NewMatrix(2, 2)
	p := NewPlayer(m, 0.5)
	if p.MatrixState != m || p.GrowthRate != 0.5 {
		t.Error("NewPlayer did not set fields correctly")
	}
}

func TestPlayerUpdateMatrix(t *testing.T) {
	m := NewMatrix(2, 2)
	p := NewPlayer(m, 1.0)
	p.UpdateMatrix(0.0)
	if m.Data[0][0] != 1.0 {
		t.Errorf("Expected m.Data[0][0]=1.0, got %v", m.Data[0][0])
	}
	if m.Data[0][1] != 0.1 || m.Data[1][0] != 0.1 || m.Data[1][1] != 0.1 {
		t.Errorf("Other elements should be 0.1, got %v", m.Data)
	}
	p.UpdateMatrix(1.0)
	if m.Data[1][1] != 1.1 {
		t.Errorf("Expected m.Data[1][1]=1.1, got %v", m.Data[1][1])
	}
	if m.Data[0][0] != 1.1 || m.Data[0][1] != 0.2 || m.Data[1][0] != 0.2 {
		t.Errorf("Other elements should be incremented by 0.1, got %v", m.Data)
	}
}

func TestPlayerGetMatrix(t *testing.T) {
	m := NewMatrix(2, 2)
	p := NewPlayer(m, 1.0)
	if p.GetMatrix() != m {
		t.Error("GetMatrix did not return correct matrix")
	}
}
