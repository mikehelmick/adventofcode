package main

import "testing"

func TestFlipFlip(t *testing.T) {
	f := FlipFlop{
		Name:         "foo",
		On:           false,
		Destinations: []string{"a"},
	}

	if s := f.Receive(HIGH, "broadcast"); s != NONE {
		t.Errorf("flip flip emitted signal on high input")
	}

	if s := f.Receive(LOW, "a"); s != HIGH {
		t.Errorf("flip flip turned on, should send high")
	}

	if s := f.Receive(LOW, "b"); s != LOW {
		t.Errorf("flip flip turned off, should send low")
	}
}

func TestConjunction(t *testing.T) {
	c := NewConjunction("foo", []string{"bar"})
	c.AddInput("a")
	c.AddInput("b")

	if c.Receive(HIGH, "a") != HIGH {
		t.Fatalf("should be low")
	}
	if c.Receive(HIGH, "b") != LOW {
		t.Fatalf("should be high")
	}
}
