package main

import (
	"fmt"
	"github.com/pkg/errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
)


func getFiles(path string) (files []string, err error) {
	var results []string
	err = filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		results = append(results, path)
		return nil
	})
	if err != nil {
		return nil, err
	}

	for i := range results {
		files = append(files, results[i])
	}

	return files, nil
}

func compare(a, b []byte) error {
	s1 := strings.TrimSpace(string(a))
	s2 := strings.TrimSpace(string(b))
	if s1 != s2 {
		return errors.New(fmt.Sprintf("got %s, expects %s", s1, s2))
	}

	return nil
}

func TestFmt(t *testing.T) {
	inputs, err := getFiles("testdata/codeformatting/in")
	if err != nil {
		t.Fatal(err)
	}
	outputs, err := getFiles("testdata/codeformatting/out")
	if err != nil {
		t.Fatal(err)
	}

	for i := range inputs {
		input, err := ioutil.ReadFile(inputs[i])
		if err != nil {
			continue
		}

		expects, err := ioutil.ReadFile(outputs[i])
		if err != nil {
			panic(err)
		}

		got, err := gofmt(input)
		if err = compare(got, expects); err != nil {
			t.Error(err)
		}
	}
}