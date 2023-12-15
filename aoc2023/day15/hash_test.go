package main

import "testing"

func TestHash(t *testing.T) {
	h := Hasher("HASH")
	if h.Hash() != 52 {
		t.Fatalf("wrong answer")
	}
}
