package main

import (
	"bufio"
	"bytes"
	// "errors"
	"flag"
	"fmt"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/charmap"
	"os"
	"os/exec"
	"strings"
)

// check if the argument is UTF-8
func isUtf(s string) bool {
	_, _, err := encoding.UTF8Validator.Transform([]byte(s), []byte(s), true)
	return err == nil
}

// check if the argument is Russian cyrillic utf-8 chars
func isCyr(s string) bool {
	for _, c := range s {
		switch {
		case 0 <= c && c <= 0x7f:
			// ascii
		case 0x410 <= c && c <= 0x44f:
			// basic russian
		case c == 0x401 || c == 0x451:
			// yo
		default:
			return false
		}
	}
	return true
}

// the interface similar to that of encoding.Decoder and encoding.Encoder
type StringTrans interface {
	String(src string) (string, error)
}

// Apply a number of transformations to the string.
func decode(src string, tlist ...StringTrans) (string, error) {
	for _, f := range tlist {
		dst, err := f.String(src)
		if err == nil {
			src = dst
			continue
		}
		// try to strip one byte at the end
		if len(src) > 4 {
			src2 := src[0 : len(src)-1]
			dst, err = f.String(src2)
			if err == nil {
				src = dst
				continue
			}
		}
		return "", err
	}
	return src, nil
}

// Extract tags into a map.
// Only the changed tags are extracted.
func extractTags(path string) (map[string]string, error) {
	cmd := exec.Command("id3info", path)
	out, err := cmd.Output()
	res := make(map[string]string)
	if err != nil {
		return res, err
	}

	win := charmap.Windows1251.NewDecoder()
	enc := charmap.Windows1251.NewEncoder()
	iso := charmap.ISO8859_1.NewEncoder()

	combinations := [][]StringTrans{
		{win},
		{enc, iso, win},
		{iso, win},
	}

	scanner := bufio.NewScanner(bytes.NewBuffer(out))
	for scanner.Scan() {
		str := scanner.Text()
		if !strings.HasPrefix(str, "=== T") {
			continue
		}
		// tag found
		if len(str) < 7 {
			continue
		}
		words := strings.SplitN(str, ":", 2)
		if len(words) != 2 {
			continue
		}
		key := str[5:8]
		value := strings.TrimSpace(words[1])
		if isUtf(value) && isCyr(value) {
			// already normal tag
			continue
		}
		newval := value
		for _, tlist := range combinations {
			val, err := decode(newval, tlist...)
			if err != nil {
				continue
			}
			if isUtf(val) && isCyr(val) {
				newval = val
				break
			}
		}
		res[key] = newval
		// fmt.Printf("%s: [%s] [%s]\n", key, value, newval)
	}
	return res, scanner.Err()
}

func main() {
	// image := flag.String("i", "", "The mp3 file to scan")
	flag.Parse()
	if len(flag.Args()) == 0 {
		fmt.Fprintln(os.Stderr, "please specify at least one mp3")
		os.Exit(1)
	}

	for _, image := range flag.Args() {
		if tags, err := extractTags(image); err != nil {
			fmt.Fprintf(os.Stderr, "failed %s: %s\n", image, err.Error())
		} else if len(tags) > 0 {
			fmt.Printf("file: %s, tags: %v\n", image, tags)
		}
	}
}
