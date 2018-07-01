package index

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestIndexer(t *testing.T) {
	d, err := ioutil.TempDir("", "test")
	if err != nil {
		t.Fatalf("Error creating temp dir: %s", err)
	}
	defer os.RemoveAll(d)
	f := filepath.Join(d, "tmp.jpg")
	f2 := filepath.Join(d, "8b", "c36727b5aa2a78e730bfd393836b246c4d565e4dc3e4f413df26e26656bb53.jpg")
	if err := ioutil.WriteFile(f, []byte("Philosophastra Illustrans"), 0600); err != nil {
		t.Fatalf("Error creating work file: %s", err)
	}

	ix, err := New(d, "file::memory:?mode=memory&cache=shared")
	if err != nil {
		t.Fatalf("Error creating indexer: %s", err)
	}
	if err := ix.AddFile(f); err != nil {
		t.Fatalf("Error adding file: %s", err)
	}

	fi, err := os.Stat(f)
	if err != nil {
		t.Fatalf("Error stating file: %s", err)
	}
	fi2, err := os.Stat(f2)
	if err != nil {
		t.Fatalf("Error stating file: %s", err)
	}
	if !os.SameFile(fi, fi2) {
		t.Errorf("Indexed file are not the same")
	}
}

func TestIndexer_merge(t *testing.T) {
	d, err := ioutil.TempDir("", "test")
	if err != nil {
		t.Fatalf("Error creating temp dir: %s", err)
	}
	defer os.RemoveAll(d)
	f := filepath.Join(d, "tmp.jpg")
	f2 := filepath.Join(d, "8b", "c36727b5aa2a78e730bfd393836b246c4d565e4dc3e4f413df26e26656bb53.jpg")
	err = ioutil.WriteFile(f, []byte("Philosophastra Illustrans"), 0600)
	if err != nil {
		t.Fatalf("Error creating work file: %s", err)
	}
	if err := os.MkdirAll(filepath.Dir(f2), 0777); err != nil {
		t.Fatalf("Error creating index dir: %s", err)
	}
	if err := ioutil.WriteFile(f2, []byte("Philosophastra Illustrans"), 0600); err != nil {
		t.Fatalf("Error creating work file: %s", err)
	}

	ix, err := New(d, "file::memory:?mode=memory&cache=shared")
	if err != nil {
		t.Fatalf("Error creating indexer: %s", err)
	}
	if err := ix.AddFile(f); err != nil {
		t.Fatalf("Error adding file: %s", err)
	}

	fi, err := os.Stat(f)
	if err != nil {
		t.Fatalf("Error stating file: %s", err)
	}
	fi2, err := os.Stat(f2)
	if err != nil {
		t.Fatalf("Error stating file: %s", err)
	}
	if !os.SameFile(fi, fi2) {
		t.Errorf("Indexed file are not the same")
	}
}
