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
	"encoding/json"
	"flag"
	"io"
	"os"
	"regexp"
	"runtime"
	"strings"
	"sync/atomic"

	"github.com/miku/parallel"
	log "github.com/sirupsen/logrus"
)

var (
	numWorkers = flag.Int("w", runtime.NumCPU(), "number of workers")
	batchSize  = flag.Int("s", 100, "batch size")
	verbose    = flag.Bool("v", false, "be verbose")
)

// Decompressor keeps options for decompression. Empty include will allow all
// files.
type Decompressor struct {
	Includes []*regexp.Regexp
	Stats    struct {
		NumArchives  int64 `json:"archives"`
		NumExtracted int64 `json:"extracted"`
		NumEntries   int64 `json:"entries"`
	}
}

func (decomp *Decompressor) ResetStats() {

}

// includesFilename return true, if this decompressor is configured to
// decompress a file with a given name.
func (decomp *Decompressor) includesFilename(name string) bool {
	if len(decomp.Includes) == 0 {
		return true
	}
	for _, include := range decomp.Includes {
		if include.MatchString(name) {
			return true
		}
	}
	return false
}

// decompressFilename takes a filename of a zipfile and returns the
// decompressed content of a fileset.
func (decomp *Decompressor) decompressFilename(name []byte) ([]byte, error) {
	defer func() {
		atomic.AddInt64(&decomp.Stats.NumArchives, 1)
	}()

	filename := strings.TrimSpace(string(name))
	r, err := zip.OpenReader(filename)
	if err != nil {
		return nil, err
	}
	defer r.Close()
	var buf bytes.Buffer

	for _, f := range r.File {
		atomic.AddInt64(&decomp.Stats.NumEntries, 1)

		if !decomp.includesFilename(f.Name) {
			continue
		}
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
		atomic.AddInt64(&decomp.Stats.NumExtracted, 1)
	}
	return buf.Bytes(), nil
}

// ArrayFlags allows to store lists of flag values.
type ArrayFlags []string

// String representation.
func (f *ArrayFlags) String() string {
	return strings.Join(*f, ", ")
}

// Set appends a value.
func (f *ArrayFlags) Set(value string) error {
	*f = append(*f, value)
	return nil
}

func main() {
	var includePatterns ArrayFlags
	flag.Var(&includePatterns, "i", "include name matching regular expression (repeatable)")
	flag.Parse()

	var includes []*regexp.Regexp

	for _, p := range includePatterns {
		c, err := regexp.Compile(p)
		if err != nil {
			log.Fatal(err)
		}
		includes = append(includes, c)
	}

	decomp := &Decompressor{Includes: includes}

	p := parallel.NewProcessor(os.Stdin, os.Stdout, decomp.decompressFilename)
	p.NumWorkers = *numWorkers
	p.BatchSize = *batchSize

	if err := p.Run(); err != nil {
		log.Fatal(err)
	}
	if *verbose {
		b, err := json.Marshal(decomp.Stats)
		if err != nil {
			log.Fatal(err)
		}
		log.Println(string(b))
	}
}
