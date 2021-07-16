package parsers

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type PcombSuite struct {
	suite.Suite
}

func (suite *PcombSuite) SetupTest() {

}

func (suite *PcombSuite) TestMatch() {
	tests := []struct {
		Parser Parser
		Input  string
		Match  bool
	}{
		{Char('a'), "a", true},
		{Char('a'), "b", false},
		{String("a"), "a", true},
		{String("abc"), "abc", true},
		{String("abcd"), "abc", false},
		{String("adbc"), "abcd", false},
		{Either(Char('a'), Char('b')), "a", true},
		{Either(Char('a'), Char('b')), "b", true},
		{Either(Char('a'), Char('b')), "x", false},
		{Either(Char('a'), String("xy")), "xyx", true},
		{Seq(Char('a'), String("b")), "ab", true},
		{Seq(Char('a'), String("b")), "ba", false},
		{Seq(Char('b'), String("a")), "ba", true},
		{Seq(Char('b'), String("a")), "ab", false},
	}

	for i, test := range tests {

		s := NewState([]byte(test.Input))
		res := s.Match(test.Parser)

		assert.Equal(suite.T(), test.Match, res.Ok, fmt.Sprintf("test [%d]. Expected %v, got %v", i, test.Match, res.Ok))
	}
}

func TestPcomb(t *testing.T) {
	suite.Run(t, new(PcombSuite))
}
