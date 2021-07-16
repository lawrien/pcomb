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

	s := NewState()

	for i, test := range tests {
		res := s.Match([]byte(test.Input), test.Parser)

		assert.Equal(suite.T(), test.Match, res.Ok, fmt.Sprintf("test [%d]. Expected %v, got %v", i, test.Match, res.Ok))
	}
}

func (suite *PcombSuite) TestMapper() {

	test := String("abc").Map(func(r Result) Result {
		if r.Ok {
			r.Bytes = []byte("Success")
		} else {
			r.Bytes = []byte("Failure")
		}
		return r
	})

	s := NewState()

	r := s.Match([]byte("abc"), test)
	assert.Equal(suite.T(), true, r.Ok, fmt.Sprintf("Expected %v, got %v", true, r.Ok))
	assert.Equal(suite.T(), "Success", string(r.Bytes), fmt.Sprintf("Expected %v, got %v", "Success", string(r.Bytes)))

	r = s.Match([]byte("cba"), test)
	assert.Equal(suite.T(), false, r.Ok, fmt.Sprintf("Expected %v, got %v", false, r.Ok))
	assert.Equal(suite.T(), "Failure", string(r.Bytes), fmt.Sprintf("Expected %v, got %v", "Failure", string(r.Bytes)))

}

func TestPcomb(t *testing.T) {
	suite.Run(t, new(PcombSuite))
}
