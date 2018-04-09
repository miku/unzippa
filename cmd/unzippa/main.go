package main

import (
	"archive/zip"
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/miku/unzippa"
)

var (
	membersFile = flag.String("m", "", "file to members to extract, one per line")
	verbose     = flag.Bool("v", false, "verbose output")
	version     = flag.Bool("version", false, "show version")
	outputFile  = flag.String("o", "", "output filename")
)

func main() {
	flag.Parse()

	if *version {
		fmt.Println(unzippa.Version)
		os.Exit(0)
	}

	if flag.NArg() == 0 {
		log.Fatal("zip file required")
	}

	if *membersFile == "" {
		log.Fatal("specify a file with members to extract via -m")
	}

	started := time.Now()

	members, err := unzippa.ReadLinesToSet(*membersFile)
	if err != nil {
		log.Fatal(err)
	}

	if *verbose {
		log.Printf("marked %d filenames for extraction", len(members))
	}

	r, err := zip.OpenReader(flag.Arg(0))
	if err != nil {
		log.Fatal(err)
	}
	defer r.Close()

	var bw *bufio.Writer

	if *outputFile != "" {
		output, err := os.Create(*outputFile)
		if err != nil {
			log.Fatal(err)
		}
		defer output.Close()
		bw = bufio.NewWriter(output)
	} else {
		bw = bufio.NewWriter(os.Stdout)
	}

	defer bw.Flush()

	var totalBytes, misses int64

	for _, f := range r.File {
		_, ok := members[f.Name]
		if !ok {
			misses++
			continue
		}
		rc, err := f.Open()
		if err != nil {
			log.Fatal(err)
		}
		n, err := io.Copy(bw, rc)
		if err != nil {
			log.Fatal(err)
		}
		rc.Close()
		totalBytes += n
	}

	if *verbose {
		log.Printf("extracted %d bytes from %d/%d files in %s",
			totalBytes, len(members), int64(len(members))+misses, time.Since(started))
	}
}
