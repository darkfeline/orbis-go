package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

type fakeIndexer struct {
	files []string
}

func (ix *fakeIndexer) AddFile(p string) error {
	ix.files = append(ix.files, p)
	return nil
}

func TestIndexFileOrDir_file(t *testing.T) {
	f, err := ioutil.TempFile("", "test")
	if err != nil {
		t.Fatalf("Error creating working file: %s", err)
	}
	defer os.Remove(f.Name())
	ix := new(fakeIndexer)
	if err := indexFileOrDir(ix, f.Name()); err != nil {
		t.Fatalf("indexFileOrDir returned error: %s", err)
	}
	exp := []string{f.Name()}
	if !reflect.DeepEqual(ix.files, exp) {
		t.Errorf("Expected %#v, got %#v", exp, ix.files)
	}
}

func TestIndexFileOrDir_dir(t *testing.T) {
	d, err := ioutil.TempDir("", "test")
	if err != nil {
		t.Fatalf("Error creating temp dir: %s", err)
	}
	defer os.RemoveAll(d)
	f1 := filepath.Join(d, "foo")
	f2 := filepath.Join(d, "bar")
	if err := ioutil.WriteFile(f1, []byte{}, 0600); err != nil {
		t.Fatalf("Error writing work file: %s", err)
	}
	if err := ioutil.WriteFile(f2, []byte{}, 0600); err != nil {
		t.Fatalf("Error writing work file: %s", err)
	}

	ix := new(fakeIndexer)
	if err := indexFileOrDir(ix, d); err != nil {
		t.Fatalf("indexFileOrDir returned error: %s", err)
	}
	exp := []string{f2, f1}
	if !reflect.DeepEqual(ix.files, exp) {
		t.Errorf("Expected %#v, got %#v", exp, ix.files)
	}
}

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
		t.Errorf("Expected %#v, got %#v", exp, got)
	}
}
