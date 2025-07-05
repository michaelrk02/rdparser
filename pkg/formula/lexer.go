package formula

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/michaelrk02/rdparser"
	"github.com/michaelrk02/rdparser/pkg/formula/pattern"
	"github.com/michaelrk02/rdparser/pkg/formula/token"
	"github.com/shivamMg/rd"
)

type Lexer struct {
	LanguagePattern *regexp.Regexp
	TokenPattern    *regexp.Regexp
}

func NewLexer() rdparser.Lexer {
	lexPattern := []string{}

	tokenDict := token.Dict()
	for _, tok := range tokenDict {
		lexPattern = append(lexPattern, tok.Hex())
	}

	lexPattern = append(lexPattern, pattern.Dict()...)

	tokenPattern := strings.Join(lexPattern, "|")
	languagePattern := fmt.Sprintf(`^\s*(\s*|%s)*\s*$`, tokenPattern)

	return &Lexer{
		LanguagePattern: regexp.MustCompile(languagePattern),
		TokenPattern:    regexp.MustCompile(tokenPattern),
	}
}

func (t *Lexer) Lex(input string) ([]rd.Token, error) {
	if !t.LanguagePattern.MatchString(input) {
		return nil, rdparser.NewLexicalError("input string is not recognizable")
	}

	tokenStrings := t.TokenPattern.FindAllString(input, -1)
	tokenResult := make([]rd.Token, len(tokenStrings))
	for i, token := range tokenStrings {
		tokenResult[i] = rdparser.Terminal(strings.ToLower(token))
	}

	return tokenResult, nil
}
