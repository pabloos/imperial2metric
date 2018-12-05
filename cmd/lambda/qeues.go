package main

import (
	"archive/zip"
	"bytes"
	"fmt"
	"imperial2metric/pkg"
	"io"
	"io/ioutil"
	"log"
	"path/filepath"
	"sync"
)

type FileMem struct {
	Name string
	Body io.Reader
}

func unzipProducer(reader io.Reader) <-chan FileMem { // 1st producer
	c := make(chan FileMem)

	body, err := ioutil.ReadAll(reader)
	if err != nil {
		fmt.Printf("error reading the zip file: %v", err)
	}

	readerAt := bytes.NewReader(body) //

	zipReader, err := zip.NewReader(readerAt, int64(len(body)))
	if err != nil {
		fmt.Println(err)
	}

	go func() {
		for _, file := range zipReader.File {
			fmt.Printf("Descomprimiendo zip: %s\n", file.Name)

			if file.Name == "es-ES/" {
				fmt.Println("encontrado el directorio")
				continue
			}

			unzippedFileBytes, err := file.Open()
			if err != nil {
				fmt.Println(err)
				continue
			}

			defer unzippedFileBytes.Close()

			c <- FileMem{filepath.Base(file.Name), unzippedFileBytes}
		}
		close(c)
	}()

	return c
}

func transformerProducer(elorigen io.Reader) <-chan FileMem { // 2nd producer
	c := make(chan FileMem)

	go func() {
		for file := range unzipProducer(elorigen) {
			fmt.Printf("Transformando %s\n", file.Name)

			readcloser := ioutil.NopCloser(file.Body) // https://stackoverflow.com/questions/52076747/how-do-i-turn-an-io-reader-into-a-io-readcloser

			result, err := pkg.TransformFile(readcloser)

			if err != nil {
				fmt.Println(fmt.Sprintf("This file %s %v", file.Name, err))
				c <- FileMem{file.Name, readcloser}
			} else {
				readcloser2 := ioutil.NopCloser(result)

				c <- FileMem{file.Name, readcloser2}
			}
		}
		close(c)
	}()

	return c
}

func zipProducer(elorigen io.Reader) io.Reader { // consumer
	buf := new(bytes.Buffer)
	w := zip.NewWriter(buf) // Create a new zip archive.

	var wg sync.WaitGroup

	wg.Add(1)

	go func() {
		for file := range transformerProducer(elorigen) {
			fmt.Printf("Almacenando zip %s\n", file.Name)

			f, err := w.Create(file.Name)
			if err != nil {
				log.Fatal(err)
			}

			buff := new(bytes.Buffer)
			buff.ReadFrom(file.Body)

			_, err = f.Write(buff.Bytes())
			if err != nil {
				log.Fatal(err)
			}
		}
		wg.Done()
	}()

	wg.Wait()

	err := w.Close()
	if err != nil {
		log.Fatal(err)
	}

	return buf
}
