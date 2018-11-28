package main

import (
	"fmt"
	"path/filepath"

	"github.com/mholt/archiver"
)

const tempDir = "/tmp/"

func isAZip(filename string) bool {
	return ".zip" == filepath.Ext(filename)
}

func decompressFile(source string) {
	archiver.Unarchive(
		fmt.Sprintf("%s%s", tempDir, source),
		fmt.Sprintf("%s%s", tempDir, "."),
	)
}
