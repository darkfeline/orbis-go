package main

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"go.felesatra.moe/orbis/internal/index"
	"go.felesatra.moe/subcommands"
	"go.felesatra.moe/xdg"
)

func init() {
	commands = append(commands, subcommands.New("index", errCmd(indexMain)))
}

func indexMain(args []string) error {
	args = args[1:]
	if len(args) == 0 {
		return nil
	}
	indexDir, err := findIndexDir(args[0])
	if err != nil {
		return errors.Wrap(err, "find index dir")
	}
	log.Printf("Found index dir %s", indexDir)
	if err := os.MkdirAll(filepath.Dir(cachePath), 0777); err != nil {
		return err
	}
	ix, err := index.New(indexDir, cachePath)
	if err != nil {
		return err
	}
	for _, f := range args {
		if err := indexFileOrDir(ix, f); err != nil {
			return err
		}
	}
	return nil
}

func indexFileOrDir(ix *index.Indexer, f string) error {
	fi, err := os.Stat(f)
	if err != nil {
		return err
	}
	if fi.IsDir() {
		return indexDir(ix, f)
	} else {
		return ix.AddFile(f)
	}
}

func indexDir(ix *index.Indexer, f string) error {
	return filepath.Walk(f, func(p string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if fi.IsDir() {
			return nil
		}
		return ix.AddFile(f)
	})
}

var cachePath = filepath.Join(xdg.CacheHome(), "orbis", "hashcache.db")

const indexDirName = "index"

func findIndexDir(p string) (string, error) {
	for last := ""; last != p; last, p = p, filepath.Dir(p) {
		fi, err := os.Stat(p)
		if err != nil {
			return "", err
		}
		if !fi.IsDir() {
			continue
		}
		fis, err := ioutil.ReadDir(p)
		if err != nil {
			return "", err
		}
		for _, fi := range fis {
			if fi.Name() == indexDirName {
				return filepath.Join(p, fi.Name()), nil
			}
		}
	}
	return "", errors.New("no index directory found")
}
