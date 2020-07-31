package main

import (
	"testing"
)

func TestGetSubs(t *testing.T) {
	CorrectCids := []string{
		"UCwi3BrUqM4xStpbCyxsb3TA",
		"UCFaYLR_1aryjfB7hLrKGRaQ",
		"UC3I2GFN_F8WudD_2jUZbojA",
	}
	Cids := GetSubs()
	if len(Cids) != len(CorrectCids) {
		t.Errorf("Incorrect number of cids were fetched")
	}
	for i := 0; i < 3; i++ {
		if Cids[i] != CorrectCids[i] {
			t.Errorf("Cids does not match the expected.")
		}
	}

}
