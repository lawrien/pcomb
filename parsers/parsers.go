package parsers

const EOL = byte('\n')

type State struct {
	buf      []byte
	pos      int
	newlines []int
}

func NewState() *State {
	return &State{}
}

func (s State) Pos() int {
	return s.pos
}

func (s *State) SetPos(p int) {
	s.pos = p
}

func (s *State) Read() (b byte) {
	if s.pos < len(s.buf) {
		b = s.buf[s.pos]
		s.pos += 1
		if b == EOL {
			s.newlines = append(s.newlines, s.pos)
		}
	}

	return
}

func (s *State) Reset() {
	s.buf = []byte{}
	s.pos = 0
	s.newlines = []int{}
}

func (s *State) Match(b []byte, p Parser) Result {
	s.Reset()
	s.buf = b

	return p(s)
}

type Result struct {
	Ok    bool
	Bytes []byte
}

type Parser func(s *State) Result
type Mapper func(r Result) Result

func (p Parser) Map(m Mapper) Parser {
	return func(s *State) Result {
		return m(p(s))
	}
}

func String(str string) Parser {
	bytes := []byte(str)
	return func(s *State) Result {
		for _, r := range bytes {
			if r != s.Read() {
				return Result{Ok: false}
			}
		}
		return Result{Ok: true, Bytes: bytes}
	}
}

func Char(r rune) Parser {
	return String(string(r))
}

func Either(parsers ...Parser) Parser {
	return func(s *State) (r Result) {
		mark := s.pos
		for _, p := range parsers {
			r = p(s)
			if r.Ok {
				break
			}
			s.pos = mark
		}
		return
	}
}

func Seq(parsers ...Parser) Parser {
	return func(s *State) (r Result) {
		mark := s.pos
		for _, p := range parsers {
			r = p(s)
			if !r.Ok {
				s.pos = mark
				break
			}
		}
		return
	}
}
