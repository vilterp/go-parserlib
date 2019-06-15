package psi

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	parserlib "github.com/vilterp/go-parserlib/pkg"
)

type CompletionsRequest struct {
	Query     string
	CursorPos parserlib.Position
}

type CompletionsResponse struct {
	Errors      []*ErrorAnnotation
	ParseError  string
	Completions []*Completion
}

func (l *Language) ServeCompletions(w http.ResponseWriter, req *http.Request) {
	reqBytes, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Println("error reading body:", err)
		w.WriteHeader(500)
		return
	}

	compReq := &CompletionsRequest{}
	if err := json.Unmarshal(reqBytes, compReq); err != nil {
		log.Println("error unmarshalling body:", err)
		w.WriteHeader(500)
		return
	}

	resp := l.GetCompletions(compReq.Query, compReq.CursorPos)

	respBytes, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, "error marshalling response", http.StatusInternalServerError)
	}

	if _, err := w.Write(respBytes); err != nil {
		log.Println("error writing response", err)
	}
}

func (l *Language) GetCompletions(query string, pos parserlib.Position) *CompletionsResponse {
	traceTree, err := l.Grammar.Parse(l.Grammar.StartRule, query, 0, nil)
	ruleTree := traceTree.ToRuleTree()
	psiTree := l.Extract(ruleTree)

	resp := &CompletionsResponse{
		Completions: l.Complete(psiTree, pos),
		Errors:      l.AnnotateErrors(psiTree),
	}
	if err != nil {
		resp.ParseError = err.Error()
	}
	return resp
}
