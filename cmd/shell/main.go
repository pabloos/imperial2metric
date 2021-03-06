package main

import (
	"bytes"
	"fmt"
	i2m "imperial2metric/pkg"
	"log"
	"os"
	"path/filepath"
)

const (
	resultsDir = "results"
)

func main() {
	files := os.Args[1:]

	if len(files) < 1 {
		fmt.Print("\nYou have not introduced any file/s\n\n")
		fmt.Print("Usage: \n --> $ imperial2metric files...\n\n")

		os.Exit(1)
	}

	for _, file := range files {
		path := filepath.Dir(file)

		dir := fmt.Sprintf("%s/%s", path, resultsDir)

		errMkdir := os.Mkdir(dir, 0755)
		if errMkdir != nil {
			fmt.Println(errMkdir)
		}

		iofile, err := os.Open(file)
		if err != nil {
			log.Fatalf("readLines: %s", err)
		}

		if filepath.Ext(iofile.Name()) == ".zip" {
			result := i2m.ZipProducer(iofile, iofile.Name())

			buff := new(bytes.Buffer)
			buff.ReadFrom(result)

			bytes := buff.Bytes()

			i2m.WriteOnFile(path, resultsDir, filepath.Base(file), bytes)

			continue
		}

		metricfile, err := i2m.TransformFile(iofile)
		if err != nil {
			fmt.Printf("Error transforming the file %s: %v", file, err)
			continue
		}

		buff := new(bytes.Buffer)
		buff.ReadFrom(metricfile)

		bytes := buff.Bytes()

		i2m.WriteOnFile(path, resultsDir, filepath.Base(file), bytes)
	}
}
