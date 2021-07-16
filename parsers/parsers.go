package parsers

import (
	"unicode/utf8"
)

const EOF = rune(0)
const EOL = rune('\n')

type State struct {
	buf      []byte
	pos      int
	newlines []int
}

func NewState(b []byte) *State {
	return &State{buf: b}
}

func (s State) Pos() int {
	return s.pos
}

func (s *State) SetPos(p int) {
	s.pos = p
}

func (s *State) Read() (r rune) {
	var size int

	if s.pos < len(s.buf) {
		r, size = utf8.DecodeRune(s.buf[s.pos:])
		s.pos += size
		if r == EOL {
			s.newlines = append(s.newlines, s.pos)
		}
	} else {
		r = EOF
	}

	return
}

func (s *State) Match(p Parser) bool {
	return p(s)
}

type Result interface {
}

type Failure struct {
	Result
}

type Success struct {
	Result
}

type Parser func(s *State) bool

func (p Parser) Map(fn func()) Parser {
	return func(s *State) bool {
		r := p(s)
		fn()
		return r
	}
}

func String(str string) Parser {
	return func(s *State) bool {
		m := s.Pos()
		for _, r := range str {
			if r != s.Read() {
				s.SetPos(m)
				return false
			}
		}
		return true
	}
}

func Char(r rune) Parser {
	return String(string(r))
}

func Either(parsers ...Parser) Parser {
	return func(s *State) (b bool) {
		for _, p := range parsers {
			b = p(s)
			if b {
				break
			}
		}
		return
	}
}

func Seq(parsers ...Parser) Parser {
	return func(s *State) (b bool) {
		m := s.Pos()
		for _, p := range parsers {
			b = p(s)
			if !b {
				s.SetPos(m)
				break
			}
		}
		return
	}
}
