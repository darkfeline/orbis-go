package hashcache

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestHashCache_missing_value(t *testing.T) {
	f, err := ioutil.TempFile("", "test")
	defer os.Remove(f.Name())
	c, err := New("file::memory:?mode=memory&cache=shared")
	if err != nil {
		t.Fatalf("Error opening cache: %s", err)
	}
	i, err := os.Stat(f.Name())
	if err != nil {
		t.Fatalf("Error stating work file: %s", err)
	}
	_, err = c.Get(f.Name(), i)
	if !IsNoRow(err) {
		t.Errorf("Expected NoRow error, got %s", err)
	}
}

type tData struct {
	data   []byte
	digest string
}

var testData = map[string]tData{
	"sophie": {
		[]byte("sophie"),
		"5e0176c9d2070a5a2a22bf74b4abed303654690d58d64221ccbd022af827abc4",
	},
	"plachta": {
		[]byte("plachta"),
		"1586e1f4475428a0424f04ec4e692376ce1c651e1ce1d9bcd280b965c970bc5e",
	},
}

func TestHashCache_get_set(t *testing.T) {
	f, err := ioutil.TempFile("", "test")
	defer os.Remove(f.Name())

	d := testData["sophie"]
	err = ioutil.WriteFile(f.Name(), d.data, 0666)
	if err != nil {
		t.Fatalf("Error writing work file: %s", err)
	}

	c, err := New("file::memory:?mode=memory&cache=shared")
	if err != nil {
		t.Fatalf("Error opening cache: %s", err)
	}
	i, err := f.Stat()
	if err != nil {
		t.Fatalf("Error stating work file: %s", err)
	}
	err = c.Set(f.Name(), i, d.digest)
	if err != nil {
		t.Fatalf("Error setting cache: %s", err)
	}
	h, err := c.Get(f.Name(), i)
	if err != nil {
		t.Fatalf("Error getting row, got %s", err)
	}
	if h != d.digest {
		t.Errorf("Expected %s, got %s", d.digest, h)
	}
}

func TestHashCache_get_changed(t *testing.T) {
	f, err := ioutil.TempFile("", "test")
	defer os.Remove(f.Name())

	d := testData["sophie"]
	err = ioutil.WriteFile(f.Name(), d.data, 0666)
	if err != nil {
		t.Fatalf("Error writing work file: %s", err)
	}

	c, err := New("file::memory:?mode=memory&cache=shared")
	if err != nil {
		t.Fatalf("Error opening cache: %s", err)
	}
	d = testData["sophie"]
	err = ioutil.WriteFile(f.Name(), d.data, 0666)
	if err != nil {
		t.Fatalf("Error writing work file: %s", err)
	}
	i, err := f.Stat()
	if err != nil {
		t.Fatalf("Error stating work file: %s", err)
	}
	err = c.Set(f.Name(), i, d.digest)
	if err != nil {
		t.Fatalf("Error setting cache: %s", err)
	}
	d = testData["plachta"]
	err = ioutil.WriteFile(f.Name(), d.data, 0666)
	if err != nil {
		t.Fatalf("Error writing work file: %s", err)
	}
	i, err = f.Stat()
	if err != nil {
		t.Fatalf("Error stating work file: %s", err)
	}
	_, err = c.Get(f.Name(), i)
	if !IsNoRow(err) {
		t.Errorf("Expected NoRow error, got %s", err)
	}
}
