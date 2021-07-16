package parsers

const EOL = byte('\n')

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

func (s *State) Match(p Parser) Result {
	return p(s)
}

type Result struct {
	Ok    bool
	Bytes []byte
}

type Parser func(s *State) Result

func (p Parser) Map(fn func()) Parser {
	return func(s *State) Result {
		r := p(s)
		fn()
		return r
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
