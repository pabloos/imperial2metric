package main

import (
	"fmt"
	i2m "imperial2metric/pkg"
	"log"
	"os"
	"path/filepath"
)

func main() {
	files := os.Args[1:]

	if len(files) < 1 {
		fmt.Print("\nYou have not introduced any file/s\n\n")
		fmt.Print("Usage: \n --> $ imperial2metric files...\n\n")

		os.Exit(1)
	}

	for _, file := range files {
		iofile, err := os.Open(file)
		if err != nil {
			log.Fatalf("readLines: %s", err)
		}

		lines := i2m.TransformFile(iofile, file)

		path := filepath.Dir(file)

		dir := fmt.Sprintf("%s/%s", path, resultsDir)

		errMkdir := os.Mkdir(dir, 0755)
		if errMkdir != nil {
			fmt.Println(errMkdir)
		}

		writeOnFile(path, resultsDir, filepath.Base(file), lines)
	}
}
