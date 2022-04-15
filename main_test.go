package main

import "testing"

func TestFoo(t *testing.T) {
	want := 1
	got := Foo()

	if got != want {
		t.Errorf("got %d, want %d", got, want)
	}
}
