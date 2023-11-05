package linecount

import (
	"bufio"
	"io"
	"log"
	"os"
)

type Token int

const (
	EOF = iota
	ILLEGAL
	NEWLINE
	CHAR
	SPACE
)

var tokens = []string{
	EOF:     "EOF",
	ILLEGAL: "ILLEGAL",
	NEWLINE: "\\n",
	SPACE:   " ",
}

func (t Token) String() string {
	return tokens[t]
}

type Position struct {
	line   int
	column int
}

type Lexer struct {
	pos    Position
	reader *bufio.Reader
}

func NewLexer(reader io.Reader) *Lexer {
	return &Lexer{
		reader: bufio.NewReader(reader),
		pos:    Position{line: 1, column: 0},
	}
}

func (l *Lexer) Lex() (Position, Token, string) {
	for {
		b, err := l.reader.ReadByte()
		if err != nil {
			if err == io.EOF {
				// We need to understand if a newline follows EOF
				// If not, then we have to manually increment the rows to get
				// the correct linecount
				if err := l.reader.UnreadByte(); err != nil {
					log.Printf("error unreading byte %s\n", err.Error())
					os.Exit(1)
				}
				// Re-read the penultimate byte to determine if it is
				// a newline
				b, err := l.reader.ReadByte()
				if err != nil {
					log.Printf("error reading byte %s\n", err.Error())
					os.Exit(1)
				}
				// If it is not, then adjust the position because it is still a line/
				// We also only run this if there are other characters on the line
				// as inidicated by a col > 0. Technically this should be taken
				// care of in the double newline check anyway.
				if b != '\n' && l.pos.column > 0 {
					l.resetPosition()
				}
				return l.pos, EOF, ""
			}
			log.Printf("error reading byte %s\n", err.Error())
			os.Exit(1)
		}
		l.pos.column++
		switch b {
		case '\n':
			// Double new line check
			// Newline on first column is illegal
			if l.pos.column == 1 {
				return l.pos, ILLEGAL, "\\n"
			}
			l.resetPosition()
			//return l.pos, NEWLINE, string(b)
		case '\r':
			ahead, err := l.reader.Peek(1)
			if err != nil {
				log.Printf("error reading ahead 1 byte %s\n", err.Error())
				os.Exit(1)
			}
			// Check if the sequence is CRLF \r\n and if so, continue
			if ahead[0] == '\n' {
				if l.pos.column == 1 { // Check if there are prior characters read
					l.pos.column = 0
				}
				continue
			} else {
				// if not throw an error
				return l.pos, ILLEGAL, "\\r"
			}
		default:
			// Check if only valid DNS characters can be used.
			// This can be replaced with another function depending
			// on what the valid characters are.
			if isValidDNSChar(b) {
				continue
				//return l.pos, CHAR, string(b)
			} else {
				if b == '\n' {
					return l.pos, ILLEGAL, "\\n"
				}
				return l.pos, ILLEGAL, string(b)
			}
		}

	}
}

func (l *Lexer) resetPosition() {
	l.pos.line++
	l.pos.column = 0
}

// Check for DNS specific characters
func isValidDNSChar(c byte) bool {
	if (c > 47 && c < 58) || (c > 64 && c < 91) || (c > 96 && c < 123) {
		return true
	}
	// I've included "." even though it is not in the RFC
	// but can still be sent through to resolve
	if c == 45 || c == 46 {
		return true
	}
	return false
}
