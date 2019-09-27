package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"
)

// Practice 1.1
func TestEcho1(t *testing.T) {
	s, sep := "", ""
	for _, arg := range os.Args {
		s += sep + arg
		sep = " "
	}
	fmt.Println(s)
}

// Practice 1.2
func TestEcho2(t *testing.T) {
	s, sep := "", ""
	for i, arg := range os.Args {
		s += sep + strconv.Itoa(i) + arg
		sep = " "
	}
	fmt.Println(s)
}

// Practice 1.3 - 1
func BenchmarkEchoNaiveJoin(b *testing.B) {
	for i := 0; i < b.N; i++ {
		s, sep := "", ""
		for i, arg := range os.Args[1:] {
			s += sep + strconv.Itoa(i) + arg
			sep = " "
		}
	}
}

// Practice 1.3 - 2
func BenchmarkEchoJoin(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = strings.Join(os.Args[1:], " ")
	}
}
