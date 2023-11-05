package linecount

import (
	"fmt"
	"log"
	"os"
	"strings"
	"testing"
)

func TestLexer(t *testing.T) {
	dir := fmt.Sprintf("%s/Documents/sec/SecLists/Discovery/DNS/", os.Getenv("HOME"))
	entries, err := os.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}
	errCount := 0
	for _, e := range entries {
		name := e.Name()
		if strings.HasSuffix(name, "txt") {
			fullname := fmt.Sprintf("%s%s", dir, e.Name())
			LineCounter, err := NewLineCounterFromFile(fullname)
			if err != nil {
				panic(err)
			}
			count, err := LineCounter.Count()
			if err != nil {
				log.Printf("\033[32mfilename:\033[0m %s \033[31merror:\033[0m %s", name, err)
				errCount++
			}
			if count > 0 {
				log.Printf("\033[32mfilename:\033[0m %s \033[32mlinecount:\033[0m %d\n", name, count)
			}
		}
	}
	log.Printf("files processed: %d errors: %d\n", len(entries), errCount)

}

func TestEdgeCases(t *testing.T) {
	t.Parallel()
	var tt = []struct {
		testName string
		s        string
		expected int
	}{
		{"One Line", "test", 1},
		{"3 Lines", "TestString\nTest\n1234", 3},
		{"Trailing newline", "TestString\nTest\n1234\n", 3},
		{"3 Lines cr lf", "TestString\r\nTest\r\n1234", 3},
	}
	for _, x := range tt {
		x := x // NOTE: https://github.com/golang/go/wiki/CommonMistakes#using-goroutines-on-loop-iterator-variables
		t.Run(x.testName, func(t *testing.T) {
			t.Parallel()
			lineCounter, err := NewLineCounterFromString(x.s)
			if err != nil {
				t.Fatalf("error %v\n", err)
			}
			count, err := lineCounter.Count()
			if err != nil {
				t.Fatalf("count error %v\n", err)
			}
			if count != x.expected {
				t.Fatalf("wrong line count! Got %d expected %d\n", count, x.expected)
			}
		})
	}
}
