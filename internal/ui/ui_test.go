package ui

import "testing"

func TestAddBattleLog(t *testing.T) {
	u := NewUI()
	u.AddBattleLog("test log")
	if len(u.battleLog) != 1 || u.battleLog[0] != "test log" {
		t.Errorf("AddBattleLog failed: got %v", u.battleLog)
	}
}

func TestClearBattleLog(t *testing.T) {
	u := NewUI()
	u.AddBattleLog("log1")
	u.AddBattleLog("log2")
	u.ClearBattleLog()
	if len(u.battleLog) != 0 {
		t.Errorf("ClearBattleLog failed: got %v", u.battleLog)
	}
}
