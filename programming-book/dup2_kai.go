package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	pos := make(map[string][]string)
	files := os.Args[1:]

	if len(files) == 0 {
		countLines(os.Stdin, pos)
	} else {
		for _, p := range files {
			f, err := os.Open(p)
			if err != nil {
				panic(err)
			}
			countLines(f, pos)
			f.Close()
		}
	}
	for line, posisions := range pos {
		if len(posisions) > 1 {
			fmt.Printf("%d\t%s\n", len(posisions), line)
			for _, p := range posisions {
				fmt.Println(p)
			}
		}
	}
}

// Practice 1.4
func countLines(f *os.File, pos map[string][]string) {
	input := bufio.NewScanner(f)
	var n int
	for input.Scan() {
		n++
		t := input.Text()
		pos[t] = append(pos[t], f.Name() + ":" + string(n))
	}
}
