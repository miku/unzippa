package main

import (
	"bufio"
	"flag"
	"io"
	"log"
	"os"
	"strings"
)

var membersFile = flag.String("m", "", "file to members to extract, one per line")

func main() {
	flag.Parse()

	if flag.NArg() == 0 {
		log.Fatal("zip file required")
	}

	members := make(map[string]bool)

	f, err := os.Open(*membersFile)
	if err != nil {
		log.Fatal(err)
	}
	br := bufio.NewReader(f)

	for {
		line, err := br.ReadString('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		line = strings.TrimSpace(line)
		members[line] = true
	}
	if err := f.Close(); err != nil {
		log.Fatal(err)
	}
	log.Printf("extracting %d members", len(members))

	// r, err := zip.OpenReader(flag.Arg(0))
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer r.Close()

	// for _, f := range r.File {}
}
