package linecount

import (
	"bytes"
	"fmt"
	"os"
)

type LineCounter struct {
	lexer     *Lexer
	lineCount int
}

func NewLineCounterFromFile(filename string) (*LineCounter, error) {
	fs, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	lexer := NewLexer(fs)
	return &LineCounter{
		lexer:     lexer,
		lineCount: 0,
	}, nil
}

func NewLineCounterFromString(lines string) (*LineCounter, error) {
	fs := bytes.NewBufferString(lines)
	lexer := NewLexer(fs)
	return &LineCounter{
		lexer:     lexer,
		lineCount: 0,
	}, nil
}

func (lc *LineCounter) Count() (int, error) {
	for {
		pos, tok, lit := lc.lexer.Lex()
		if tok == EOF {
			lc.lineCount = pos.line - 1
			break
		}
		if tok == ILLEGAL {
			return -1, fmt.Errorf("invalid character '%s' found at row %d col %d\n", lit, pos.line, pos.column)
		}
	}
	return lc.lineCount, nil
}
