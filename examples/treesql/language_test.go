package treesql

import (
	"testing"

	parserlib "github.com/vilterp/go-parserlib/pkg"
)

func TestLanguageServer(t *testing.T) {
	l := MakeLanguage(BlogSchema)

	c := l.GetCompletions(`MANY `, parserlib.Position{Line: 1, Col: 6})
	t.Logf(`%+v`, c)
}
