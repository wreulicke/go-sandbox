package javap

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/tinycedar/classp/classfile"
)

func TestParseClassFile(t *testing.T) {
	bs, err := ioutil.ReadFile("testdata/Test.class")
	if err != nil {
		t.Fatal(err)
	}
	cp := classfile.Parse(bs)
	for _, m := range cp.Methods() {
		fmt.Println(m.Name())
	}
}
