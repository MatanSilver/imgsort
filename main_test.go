package main

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestCopy(t *testing.T) {
	out, err := os.Create("./testfile.txt")
	defer out.Close()
	if err != nil {
		panic(err)
	}
	_, err = out.WriteString("this is a test\n")
	out.Sync()

	Copy("./testfile.txt", "./newtestfile.txt")

	dat, err := ioutil.ReadFile("./newtestfile.txt")
	if err != nil {
		panic(err)
	}
	recovered_string := string(dat)
	if "this is a test\n" != recovered_string {
		t.Errorf("Recovered string doesn't match written string")
	}

	err = os.Remove("./testfile.txt")
	if err != nil {
		panic(err)
	}
	err = os.Remove("./newtestfile.txt")
	if err != nil {
		panic(err)
	}
}
