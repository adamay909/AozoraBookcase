package main

import (
	"io"
	"log"
	"os"
	"path/filepath"
)

func rootIsSafe() bool {

	if filepath.IsAbs(root) {

		return false

	}

	home, _ := os.LookupEnv("HOME")
	os.Chdir(home)

	if filepath.IsLocal(root) {
		return true
	}

	return false
}

func setupLogging() {
	var w io.Writer

	f, err := os.OpenFile(filepath.Join(root, "aozora.log"), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)

	if err != nil {
		log.Println(err)
	}

	if verbose {
		w = io.MultiWriter(f, os.Stdout)
	} else {
		w = f
	}
	log.SetOutput(w)
	return
}
