package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"go.felesatra.moe/subcommands"
)

func init() {
	commands = append(commands, subcommands.New("bucket", errCmd(bucketMain)))
}

func bucketMain(args []string) error {
	fs := flag.NewFlagSet(args[0], flag.ExitOnError)
	var rootDir string
	fs.StringVar(&rootDir, "root", ".", "Bucket root directory")
	_ = fs.Parse(args[1:])
	var files []string
	if fs.NArg() == 0 {
		var err error
		files, err = unbucketedFiles(rootDir)
		if err != nil {
			return err
		}
	} else {
		files = fs.Args()
	}
	buckets, err := getBuckets(rootDir)
	if err != nil {
		return err
	}
	for _, f := range files {
		for _, b := range buckets {
			if strings.Contains(f, b) {
				log.Printf("Moving %s into %s", f, b)
				os.Rename(f, filepath.Join(rootDir, b, filepath.Base(f)))
				break
			}
		}
	}
	return nil
}

func getBuckets(p string) ([]string, error) {
	fis, err := ioutil.ReadDir(p)
	if err != nil {
		return nil, err
	}
	var fs []string
	for _, fi := range fis {
		if fi.Mode()&os.ModeSymlink != 0 {
			fi, err = os.Stat(filepath.Join(p, fi.Name()))
			if err != nil {
				log.Printf("get buckets: skipping symlink %s: %s", fi.Name(), err)
			}
		}
		if fi.IsDir() {
			fs = append(fs, fi.Name())
		}
	}
	return fs, nil
}

func unbucketedFiles(p string) ([]string, error) {
	fis, err := ioutil.ReadDir(p)
	if err != nil {
		return nil, err
	}
	var fs []string
	for _, fi := range fis {
		if fi.Mode().IsRegular() {
			fs = append(fs, fi.Name())
		}
	}
	return fs, nil
}
