package formula

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/michaelrk02/rdparser"
	"github.com/michaelrk02/rdparser/pkg/formula/token"
	"github.com/shivamMg/rd"
)

type Tokenizer struct {
	R *regexp.Regexp
}

func NewTokenizer() *Tokenizer {
	dict := token.Dict()

	dictPatternArr := make([]string, len(dict))
	for i, tok := range dict {
		dictPatternArr[i] = tok.Hex()
	}
	dictPattern := strings.Join(dictPatternArr, "|")

	return &Tokenizer{
		R: regexp.MustCompile(fmt.Sprintf(`[a-zA-Z][a-zA-Z0-9]*|[0-9.]+|\[[a-zA-Z0-9-:]+\]|%s`, dictPattern)),
	}
}

func (t *Tokenizer) Tokenize(input string) ([]rd.Token, error) {
	tokenStrings := t.R.FindAllString(input, -1)
	tokenResult := make([]rd.Token, len(tokenStrings))
	for i, token := range tokenStrings {
		tokenResult[i] = rdparser.Terminal(strings.ToLower(token))
	}
	return tokenResult, nil
}
