package ui

import "testing"

func TestAddBattleLog(t *testing.T) {
	tests := []struct {
		name    string
		logs    []string
		wantLen int
	}{
		{"single log", []string{"test log"}, 1},
		{"multiple logs", []string{"a", "b", "c"}, 3},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := NewUI()
			for _, log := range tt.logs {
				u.AddBattleLog(log)
			}
			if len(u.battleLog) != tt.wantLen {
				t.Errorf("AddBattleLog failed: got %v", u.battleLog)
			}
		})
	}
}

func TestClearBattleLog(t *testing.T) {
	tests := []struct {
		name string
		logs []string
	}{
		{"clear after add", []string{"log1", "log2"}},
		{"clear empty", []string{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := NewUI()
			for _, log := range tt.logs {
				u.AddBattleLog(log)
			}
			u.ClearBattleLog()
			if len(u.battleLog) != 0 {
				t.Errorf("ClearBattleLog failed: got %v", u.battleLog)
			}
		})
	}
}
