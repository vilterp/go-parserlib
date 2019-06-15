package treesql

import (
	"testing"

	parserlib "github.com/vilterp/go-parserlib/pkg"
)

func TestLanguageServer(t *testing.T) {
	l := MakeLanguage(BlogSchema)

	c := l.GetCompletions(``, parserlib.Position{Line: 1, Col: 1, Offset: 0})
	t.Log(c)
}
