package main

import (
	"archive/zip"
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/miku/unzippa"
)

var (
	membersFile = flag.String("m", "", "file to members to extract, one per line")
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

	members, err := unzippa.ReadLinesToSet(*membersFile)
	if err != nil {
		log.Fatal(err)
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
		defer bw.Flush()
	} else {
		bw = bufio.NewWriter(os.Stdout)
		defer bw.Flush()
	}

	for _, f := range r.File {
		_, ok := members[f.Name]
		if !ok {
			continue
		}
		rc, err := f.Open()
		if err != nil {
			log.Fatal(err)
		}
		_, err = io.Copy(bw, rc)
		if err != nil {
			log.Fatal(err)
		}
		rc.Close()
	}
}
