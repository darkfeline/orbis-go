// Package index implements indexing and merging files in a directory
package index

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"go.felesatra.moe/orbis/internal/hashcache"
)

// Indexer provides file indexing methods.
type Indexer struct {
	indexDir string
	cache    *hashcache.HashCache
}

// New creates a new Indexer.
func New(indexDir, cachePath string) (*Indexer, error) {
	c, err := hashcache.New(cachePath)
	if err != nil {
		return nil, err
	}
	return &Indexer{
		indexDir: indexDir,
		cache:    c,
	}, nil
}

// AddFile adds a file to an index.  If there is already a different
// file in the index with the same hash value, the error will satisfy
// IsCollision.
func (ix *Indexer) AddFile(path string) error {
	digest, err := ix.fileDigest(path)
	if err != nil {
		return err
	}
	return ix.internFile(path, digest)
}

func (ix *Indexer) fileDigest(path string) (string, error) {
	fi, err := os.Stat(path)
	if err != nil {
		return "", err
	}
	digest, err := ix.cache.Get(path, fi)
	if err != nil {
		if !hashcache.IsNoRow(err) {
			return "", err
		}
		digest, err = fileDigest(path)
		if err != nil {
			return "", err
		}
		err = ix.cache.Set(path, fi, digest)
	}
	return digest, err
}

func (ix *Indexer) internFile(path, digest string) error {
	ext := filepath.Ext(path)
	dst := filepath.Join(ix.indexDir, digest[:2], fmt.Sprintf("%s%s", digest[2:], ext))
	return mergeLink(path, dst)
}

func fileDigest(path string) (string, error) {
	d, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256(d)
	return fmt.Sprintf("%x", sum), nil
}

type collisionErr struct{}

func (e collisionErr) Error() string {
	return "file hash collision"
}

// IsCollision returns true if the error indicates a hash value
// collision.
func IsCollision(e error) bool {
	_, ok := e.(collisionErr)
	return ok
}

// mergeLink links src to dst.  If dst exists and is the same file as
// src, do nothing.  If dst exists, is a different file, and has the
// same contents, replace dst with a link to src.  If dst exists and
// has different contents, return an error for which IsCollision is
// true.
func mergeLink(src, dst string) error {
	fi, err := os.Stat(dst)
	if err != nil {
		if err := os.MkdirAll(filepath.Dir(dst), 0777); err != nil {
			return err
		}
		return os.Link(src, dst)
	}
	fi2, err := os.Stat(src)
	if err != nil {
		return err
	}
	if os.SameFile(fi, fi2) {
		return nil
	}
	d1, err := ioutil.ReadFile(src)
	if err != nil {
		return err
	}
	d2, err := ioutil.ReadFile(dst)
	if err != nil {
		return err
	}
	if !bytes.Equal(d1, d2) {
		return collisionErr{}
	}
	if err = os.Remove(dst); err != nil {
		return err
	}
	return os.Link(src, dst)
}
