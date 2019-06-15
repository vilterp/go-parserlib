package psi

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	parserlib "github.com/vilterp/go-parserlib/pkg"
)

type CompletionsRequest struct {
	Query  string
	Offset int
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

	resp := l.getCompletions(compReq.Query, compReq.Offset)

	respBytes, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, "error marshalling response", http.StatusInternalServerError)
	}

	if _, err := w.Write(respBytes); err != nil {
		log.Println("error writing response", err)
	}
}

func (l *Language) getCompletions(query string, offset int) *CompletionsResponse {
	psiTree, err := l.Parse(string(query))
	if err != nil {
		return &CompletionsResponse{
			ParseError: err.Error(),
		}
	}

	pos := parserlib.PositionFromOffset(query, offset)

	resp := &CompletionsResponse{
		Completions: l.Complete(psiTree, pos),
		Errors:      l.AnnotateErrors(psiTree),
	}

	return resp
}
