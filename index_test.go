package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestFindIndexDir(t *testing.T) {
	d, err := ioutil.TempDir("", "test")
	if err != nil {
		t.Fatalf("Error creating temp dir: %s", err)
	}
	defer os.RemoveAll(d)
	if err := os.MkdirAll(filepath.Join(d, "foo", "bar", "baz"), 0700); err != nil {
		t.Fatalf("Error creating work dirs: %s", err)
	}
	if err := os.MkdirAll(filepath.Join(d, "foo", "index"), 0700); err != nil {
		t.Fatalf("Error creating work dirs: %s", err)
	}
	got, err := findIndexDir(filepath.Join(d, "foo", "bar", "baz"))
	if err != nil {
		t.Fatalf("findIndexDir returned an error: %s", err)
	}
	exp := filepath.Join(d, "foo", "index")
	if got != exp {
		t.Errorf("Expected %s, got %s", exp, got)
	}
}
