package test_harness

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	parserlib "github.com/vilterp/go-parserlib/pkg"
	"github.com/vilterp/go-parserlib/pkg/psi"
)

type completionsRequest struct {
	Input     string
	CursorPos int // TODO: line/col?
}

type completionsResponse struct {
	TraceTree        *parserlib.TraceTree
	RuleTree         *parserlib.Node
	PSITree          psi.Node
	Completions      psi.Completions
	ErrorAnnotations []*psi.ErrorAnnotation
	Err              string
}

// TODO: use some logging middleware
// which prints statuses, urls, and times

type server struct {
	language          *psi.Language
	serializedGrammar *parserlib.SerializedGrammar

	mux *http.ServeMux
}

func NewServer(l *psi.Language) *server {
	mux := http.NewServeMux()

	// Serve UI static files.
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("/index.html")
		http.ServeFile(w, r, "pkg/test_harness/build/index.html")
	})

	fileServer := http.FileServer(http.Dir("pkg/test_harness/build/static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fileServer))

	server := &server{
		language:          l,
		serializedGrammar: l.Grammar.Serialize(),
		mux:               mux,
	}

	// Serve grammar and completions.
	mux.HandleFunc("/grammar", server.handleGrammar)
	mux.HandleFunc("/completions", server.handleCompletions)

	return server
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

func (s *server) handleGrammar(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if err := json.NewEncoder(w).Encode(s.serializedGrammar); err != nil {
		log.Println("err encoding json:", err)
	}
	end := time.Now()
	log.Println("/grammar responded in", end.Sub(start))
}

func (s *server) handleCompletions(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if r.Method != "POST" {
		log.Println("/completions: expecting GET")
		http.Error(w, "expecting GET", 400)
		return
	}
	// Decode request.
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	var cr completionsRequest
	err := decoder.Decode(&cr)
	if err != nil {
		log.Printf("/completions error: %v", err)
		http.Error(w, fmt.Sprintf("error parsing request body: %v", err), 400)
		return
	}

	resp, err := s.getResp(&cr)

	// Respond.
	if err := json.NewEncoder(w).Encode(&resp); err != nil {
		log.Println("err encoding json:", err)
		http.Error(w, err.Error(), 500)
	}

	end := time.Now()
	log.Println("/completions responded in", end.Sub(start))
}

func (s *server) getResp(req *completionsRequest) (*completionsResponse, error) {
	resp := &completionsResponse{}

	// Parse it.
	trace, err := s.language.Grammar.Parse("select", req.Input, req.CursorPos, nil)
	resp.TraceTree = trace
	if err != nil {
		resp.Err = err.Error()
	}
	if trace == nil {
		return resp, nil
	}

	// Get completions.
	syntaxCompletions, err := trace.GetCompletions()
	if err != nil {
		resp.Err = err.Error()
	}
	resp.Completions = append(resp.Completions, makeSyntaxCompletions(syntaxCompletions)...)

	// Get rule tree.
	resp.RuleTree = trace.ToTree()

	// Get PSI tree.
	resp.PSITree = s.language.Extract(resp.RuleTree)
	resp.ErrorAnnotations = s.language.AnnotateErrors(resp.PSITree)
	// TODO(vilterp): get position from cursor offset
	cursorPos := parserlib.Position{
		Line:   1,
		Col:    1,
		Offset: 0,
	}
	resp.Completions = append(
		resp.Completions,
		s.language.Complete(resp.PSITree, cursorPos)...,
	)

	return resp, nil
}

func makeSyntaxCompletions(cs []string) psi.Completions {
	var out psi.Completions
	for _, c := range cs {
		out = append(out, &psi.Completion{
			Kind:    "syntax",
			Content: c,
		})
	}
	return out
}
