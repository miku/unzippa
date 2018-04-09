package unzippa

import (
	"bufio"
	"io"
	"os"
	"strings"
)

// ReadLinesToSet reads lines of a file into a set like hash map.
func ReadLinesToSet(filename string) (map[string]bool, error) {
	result := make(map[string]bool)

	f, err := os.Open(filename)
	if err != nil {
		return result, err
	}
	defer f.Close()

	br := bufio.NewReader(f)
	for {
		line, err := br.ReadString('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			return result, err
		}
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		result[line] = true
	}
	return result, nil
}
