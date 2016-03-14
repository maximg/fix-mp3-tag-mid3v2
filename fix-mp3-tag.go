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

type StringTrans interface {
	String(src string) (string, error)
}

func showerr(err error) string {
	if err != nil {
		return err.Error()
	}
	return ""
}

func decode(src string, tlist ...StringTrans) (string, error) {
	// fmt.Println("# ---")
	for _, f := range tlist {
		dst, err := f.String(src)
		// fmt.Printf("#%d  [%s]/%d -> [%s]/%d %s\n",
		//   i, src, len(src), dst, len(dst), showerr(err))
		if err == nil {
			src = dst
			continue
		}
		// try to strip one byte at the end
		if len(src) > 4 {
			src2 := src[0 : len(src)-1]
			dst, err = f.String(src2)
			// fmt.Printf("#%d+ [%s]/%d -> [%s]/%d %s\n",
			//   i, src2, len(src2), dst, len(dst), showerr(err))
			if err == nil {
				src = dst
				continue
			}
		}
		return "", err
	}
	return src, nil
}

// try to convert the argument from win1251 to utf8
func fromWin(src string) (string, error) {
	w := charmap.Windows1251.NewDecoder()
	return w.String(src)
}

func fromLat1(src string) (string, error) {
	e := charmap.ISO8859_1.NewEncoder()
	tmp, err := e.String(src)
	if err != nil {
		return "", err
	}
	return fromWin(tmp)
}

func extractTags(path string) error {
	cmd := exec.Command("id3info", path)
	out, err := cmd.Output()
	if err != nil {
		return err
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
		oldval := value
		if !(isUtf(value) && isCyr(value)) {
			for _, tlist := range combinations {
				newval, err := decode(value, tlist...)
				if err != nil {
					continue
				}
				if isUtf(newval) && isCyr(newval) {
					value = newval
					break
				}
			}
		}
		fmt.Printf("%s: [%s] [%s]\n", key, oldval, value)
	}
	return scanner.Err()
}

func main() {
	// image := flag.String("i", "", "The mp3 file to scan")
	flag.Parse()
	if len(flag.Args()) == 0 {
		fmt.Fprintln(os.Stderr, "please specify at least one mp3")
		os.Exit(1)
	}

	for _, image := range flag.Args() {
		fmt.Println("###", image)
		if err := extractTags(image); err != nil {
			fmt.Fprintf(os.Stderr, "failed %s: %s\n", image, err.Error())
		}
	}
}
