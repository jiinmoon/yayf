package main

import (
	"testing"
)

func TestGetSubs(t *testing.T) {
	CorrectCids := []string{
		"UCwi3BrUqM4xStpbCyxsb3TA",
		"UC3I2GFN_F8WudD_2jUZbojA",
	}
	CorrectPids := []string{
		"PL1L0fRHNDxrKzTrezXURug0IWOZQy4hg7",
	}
	subscriptions := GetSubs()
	// Check lengths of Cids and Pids
	if len(subscriptions.Cids) != len(CorrectCids) {
		t.Errorf("Incorrect number of cids were fetched")
	}
	if len(subscriptions.Pids) != len(CorrectPids) {
		t.Errorf("Incorrect number of cids were fetched")
	}
	// Check ids vs Expected
	for i := 0; i < len(CorrectCids); i++ {
		if subscriptions.Cids[i] != CorrectCids[i] {
			t.Errorf("Cids does not match the expected.")
		}
	}
	for i := 0; i < len(CorrectPids); i++ {
		if subscriptions.Cids[i] != CorrectCids[i] {
			t.Errorf("Cids does not match the expected.")
		}
	}

}
