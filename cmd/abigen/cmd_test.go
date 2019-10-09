package main

import "testing"

func TestGenerate(t *testing.T) {
	err := Generate("example", "main", "token", "example/token.json", true)
	if err != nil {
		t.Error(err)
	}
}
