package pkg

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"
	"sync"
)

//FileMem represents basic file (name + content/reader)
type FileMem struct {
	Name string
	Body io.Reader
}

func unzipProducer(reader io.Reader, filename string) <-chan FileMem { // 1st producer
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
			if file.Name == fmt.Sprintf("%s/", strings.TrimSuffix(filename, filepath.Ext(filename))) {
				fmt.Println("encontrado el directorio")
				continue
			}

			unzippedFileBytes, err := file.Open()
			if err != nil {
				fmt.Println(err)
				continue
			}

			// defer unzippedFileBytes.Close() // read after close error

			c <- FileMem{filepath.Base(file.Name), unzippedFileBytes}
		}
		close(c)
	}()

	return c
}

func transformerProducer(elorigen io.Reader, filename string) <-chan FileMem { // 2nd producer
	c := make(chan FileMem)

	go func() {
		for file := range unzipProducer(elorigen, filename) {
			readcloser := ioutil.NopCloser(file.Body) // https://stackoverflow.com/questions/52076747/how-do-i-turn-an-io-reader-into-a-io-readcloser

			result, err := TransformFile(readcloser)

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

func ZipProducer(elorigen io.Reader, filename string) io.Reader { // consumer
	buf := new(bytes.Buffer)
	w := zip.NewWriter(buf) // Create a new zip archive.

	var wg sync.WaitGroup

	wg.Add(1)

	go func() {
		for file := range transformerProducer(elorigen, filename) {
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
