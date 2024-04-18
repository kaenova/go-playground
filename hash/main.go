package main

import (
	"crypto/md5"
	"fmt"
	"log"
	"os"
	"path"
)

func main() {
	// Get current working directory
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal("Cannot get current working directory")
	}

	// List file on current working directory
	dirFiles, err := os.ReadDir(cwd)
	if err != nil {
		log.Fatal("Cannot list files on current working directory")
	}

	// Create new list of files from dirFiles
	files := []string{}
	for _, file := range dirFiles {
		if !file.IsDir() {
			files = append(files, file.Name())
		}
	}

	// Hash files
	hashedFiles := map[string]string{}
	for _, file := range files {
		f, err := os.ReadFile(path.Join(cwd, file))
		if err != nil {
			log.Fatal("Cannot read file " + file)
		}
		hashedFiles[file] = fmt.Sprintf("%x", md5.Sum(f))
	}

	// Write map to current working directory
	f, err := os.Create(path.Join(cwd, "MD5.txt"))
	if err != nil {
		log.Fatal("Cannot create file MD5.txt")
	}

	for file, hash := range hashedFiles {
		f.WriteString(file + "\n" + hash + "\n\n")
	}
}