package lexer

import (
	"fmt"
	"io"
)

type Type string

const (
	CommaType      Type = "Comma"
	EOFType        Type = "EOF"
	EqualType      Type = "Equal"
	ErrorType      Type = "Error"
	IdentifierType Type = "Identifier"
	IntegerType    Type = "Integer"
	StringType     Type = "String"
	StarType       Type = "Star"
	WhitespaceType Type = "Whitespace"
)

func (t Type) String() string {
	return string(t)
}

type Lexer interface {
	Next() bool
	Token() *Token
	UnreadToken() error
}

type Token struct {
	Raw  string
	Type Type
}

type tokenizer struct {
	stream         io.RuneScanner
	current        *Token
	previous       *Token
	unread         bool
	skipWhitespace bool
}

func New(stream io.RuneScanner) Lexer {
	return &tokenizer{stream: stream}
}

func NewFilterWhitespace(stream io.RuneScanner) Lexer {
	return &tokenizer{stream: stream,
		skipWhitespace: true,
	}
}

func (l *tokenizer) Next() bool {
	if l.unread {
		l.unread = false
		return true
	}

	var token *Token
	for {
		// Lex until there is a complete token.
		for next := l.initial; next != nil; token, next = next() {
		}
		// If keeping whitespace, or the token is not whitespace, break out.
		if !l.skipWhitespace || token.Type != WhitespaceType {
			break
		}
	}

	l.previous = l.current
	l.current = token
	return true
}

type lexerFn func() (*Token, lexerFn)

func (l *tokenizer) initial() (*Token, lexerFn) {
	r, _, err := l.stream.ReadRune()
	if err == io.EOF {
		return &Token{Type: EOFType}, nil
	}
	if err != nil {
		return &Token{
			Type: ErrorType,
			Raw:  fmt.Sprintf("unable to read rune: %+v", err),
		}, nil
	}
	if r == ' ' || r == '\n' {
		l.stream.UnreadRune()
		return nil, l.whitespace
	} else if 'a' <= r && r <= 'z' || 'A' <= r && r <= 'Z' {
		l.stream.UnreadRune()
		return nil, l.identifier
	} else if '0' <= r && r <= '9' {
		l.stream.UnreadRune()
		return nil, l.number
	} else if '\'' == r {
		l.stream.UnreadRune()
		return nil, l.string
	} else if r == ',' {
		return &Token{Type: CommaType, Raw: ","}, nil
	} else if r == '=' {
		return &Token{Type: EqualType, Raw: "="}, nil
	} else if r == '*' {
		return &Token{Type: StarType, Raw: "*"}, nil
	} else {
		return &Token{
			Type: ErrorType,
			Raw:  fmt.Sprintf("unrecognized char while tokenizing: %q", r),
		}, nil
	}
}

func (l *tokenizer) identifier() (*Token, lexerFn) {
	raw := ""
	for {
		r, _, err := l.stream.ReadRune()
		if err == io.EOF {
			l.stream.UnreadRune()
			return &Token{Type: IdentifierType, Raw: raw}, nil
		}
		if err != nil {
			return &Token{Type: ErrorType}, nil
		}
		if 'a' <= r && r <= 'z' ||
			'A' <= r && r <= 'Z' ||
			'0' <= r && r <= '9' ||
			r == '_' {
			raw += string(r)
		} else {
			l.stream.UnreadRune()
			return &Token{Type: IdentifierType, Raw: raw}, nil
		}
	}
}

func (l *tokenizer) number() (*Token, lexerFn) {
	raw := ""
	for {
		r, _, err := l.stream.ReadRune()
		if err == io.EOF {
			l.stream.UnreadRune()
			return &Token{Type: IntegerType, Raw: raw}, nil
		}
		if err != nil {
			return &Token{Type: ErrorType}, nil
		}
		if '0' <= r && r <= '9' {
			raw += string(r)
		} else {
			l.stream.UnreadRune()
			return &Token{Type: IntegerType, Raw: raw}, nil
		}
	}
}

func (l *tokenizer) string() (*Token, lexerFn) {
	r, _, _ := l.stream.ReadRune()
	raw := string(r)
	for {
		r, _, err := l.stream.ReadRune()
		if err == io.EOF {
			l.stream.UnreadRune()
			return &Token{Type: IdentifierType, Raw: raw}, nil
		}
		if err != nil {
			return &Token{Type: ErrorType}, nil
		}
		raw += string(r)
		if '\'' == r {
			return &Token{Type: StringType, Raw: raw}, nil
		}
	}
}

func (l *tokenizer) whitespace() (*Token, lexerFn) {
	raw := ""
	for {
		r, _, err := l.stream.ReadRune()
		if err == io.EOF {
			l.stream.UnreadRune()
			return &Token{Type: WhitespaceType, Raw: raw}, nil
		}
		if err != nil {
			return &Token{Type: ErrorType}, nil
		}
		switch r {
		case ' ', '\n':
			raw += string(r)
		default:
			l.stream.UnreadRune()
			return &Token{Type: WhitespaceType, Raw: raw}, nil
		}
	}
}

func (l *tokenizer) Token() *Token {
	return l.current
}

func (l *tokenizer) UnreadToken() error {
	if l.unread {
		return fmt.Errorf("previous operation was not Next()")
	}
	l.unread = true
	return nil
}
