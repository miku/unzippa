// The unzippall tools takes a list of files and extracts them to stdout in
// parallel. Order is not preserved.
//
// Usage:
//
//     $ find /tmp/updates -type f -name "*zip" | unzippall > data.file
//
package main

import (
	"archive/zip"
	"bytes"
	"io"
	"log"
	"os"
	"strings"

	"github.com/miku/parallel"
)

func decompressFile(b []byte) ([]byte, error) {
	filename := strings.TrimSpace(string(b))
	r, err := zip.OpenReader(filename)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	var buf bytes.Buffer

	for _, f := range r.File {
		rc, err := f.Open()
		if err != nil {
			return nil, err
		}
		_, err = io.Copy(&buf, rc)
		if err != nil {
			return nil, err
		}
		rc.Close()
		io.WriteString(&buf, "\n")
	}
	return buf.Bytes(), nil
}

func main() {
	p := parallel.NewProcessor(os.Stdin, os.Stdout, decompressFile)
	if err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
