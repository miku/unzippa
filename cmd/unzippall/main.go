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
	"flag"
	"io"
	"log"
	"os"
	"runtime"
	"strings"

	"github.com/miku/parallel"
)

var (
	numWorkers = flag.Int("w", runtime.NumCPU(), "number of workers")
	batchSize  = flag.Int("s", 100, "batch size")
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
	flag.Parse()

	p := parallel.NewProcessor(os.Stdin, os.Stdout, decompressFile)
	p.NumWorkers = *numWorkers
	p.BatchSize = *batchSize

	if err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
