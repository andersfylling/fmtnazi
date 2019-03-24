package main

import (
	"fmt"
	"testing"
)

func TestFmt(t *testing.T) {
	code := "```go\n" +
		`func main() {
fmt.Println("test")
}` + "```"

	s, err := gofmt(code)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(s)
}