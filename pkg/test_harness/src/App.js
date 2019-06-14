import React, { Component } from 'react';
import "./App.css";
import { GrammarView } from './GrammarView';
import {formatSpan} from "./span";
import {RuleTreeView} from "./RuleTreeView";

const INITIAL_QUERY = `MANY blog_posts {
  id,
  title
}`;

// TODO: bundle all these "onhighlightX" and "highlightedX" props to a
// highlightContext or something

class App extends Component {
  constructor() {
    super();
    this.state = {
      query: INITIAL_QUERY,
      cursorPos: 0,
      grammar: null,

      resp: null,

      highlightedRuleID: null,
      highlightedSpan: null,
    };
    this.textArea = React.createRef();
  }

  componentDidMount() {
    this.fetchGrammar();
    this.fetchQuery(0, INITIAL_QUERY);
  }

  fetchGrammar() {
    fetch("/grammar").then((resp) => {
      resp.json().then((grammar) => {
        console.log("grammar:", grammar);
        this.setState({
          grammar: grammar,
        })
      }).catch((err) => {
        console.error("error parsing grammar:", err);
      });
    });
  }

  fetchQuery(pos, query) {
    const completionReq = {
      Input: query,
      CursorPos: pos,
    };
    fetch("/completions", {
      method: "POST",
      body: JSON.stringify(completionReq),
    }).then((resp) => {
      resp.json().then((completionResp) => {
        console.log("trace:", completionResp);
        this.setState({
          resp: completionResp,
        });
      }).catch((err) => {
        console.error("error parsing completions:", err);
      });
    });
  }

  handleHighlightRule = (ruleID, highlight) => {
    this.setState({
      highlightedRuleID: highlight ? ruleID : null,
    });
  }

  handleHighlightSpan = (span, highlight) => {
    this.setState({
      highlightedSpan: highlight ? span : null,
    });
  }

  updateQueryAndPos = () => {
    const textArea = this.textArea.current;
    const query = textArea.value;
    this.setState({
      cursorPos: textArea.selectionStart,
      query: query,
    });
    this.fetchQuery(textArea.selectionStart, query);
  }

  render() {
    const highlightProps = {
      onHighlightSpan: this.handleHighlightSpan,
      highlightedSpan: this.state.highlightedSpan,
      onHighlightRule: this.handleHighlightRule,
      highlightedRuleID: this.state.highlightedRuleID,
    };

    return (
      <div className="App">
        <h1>TreeSQL Parser Test Harness</h1>
        <div className="app-container">
          <div className="grid-cell app-editor">
            <textarea
              ref={this.textArea}
              style={{
                fontFamily: "monospace",
                height: "100%",
                width: "100%",
              }}
              rows={10}
              value={this.state.query}
              onChange={this.updateQueryAndPos}
              onKeyUp={this.updateQueryAndPos}
              onClick={this.updateQueryAndPos}
            />
            <br />
            Pos: {this.state.cursorPos}<br />
            {this.state.resp && this.state.resp.Completions
              ? <ul>
                  {this.state.resp.Completions.map((completion) => (
                    <li key={completion.Content}>{completion.Kind}: {completion.Content}</li>
                  ))}
                </ul>
              : null}
          </div>
          <div className="grid-cell app-sourceview">
            {/*{this.state.resp && this.state.grammar*/}
            {/*  ? <SourceView*/}
            {/*      trace={this.state.resp.TraceTree}*/}
            {/*      grammar={this.state.grammar}*/}
            {/*      {...highlightProps}*/}
            {/*    />*/}
            {/*  : <span>&lt;don't have both trace & grammar yet&gt;</span>}*/}
            {this.state.resp && this.state.grammar
              ? <>
                  <h3>Errors</h3>
                  <ul>
                    {(this.state.resp.ErrorAnnotations || []).map(err => (
                      <li><code>{formatSpan(err.Span)}: {err.Message}</code></li>
                    ))}
                  </ul>
                </>
              : <span>&lt;don't have both trace & grammar yet&gt;</span>}
          </div>
          <div className="grid-cell app-traceview">
            <h3>Rule Tree</h3>
            {this.state.resp && this.state.grammar
              ? <RuleTreeView node={this.state.resp.RuleTree} />
              : <span>&lt;don't have both trace & grammar yet&gt;</span>}
          </div>
          <div className="grid-cell app-grammar">
            <h3>Grammar</h3>
            {this.state.grammar
              ? <GrammarView
                  grammar={this.state.grammar}
                  {...highlightProps}
                />
              : <span>&lt;don't have grammar yet&gt;</span>}
          </div>
        </div>
      </div>
    );
  }
}

export default App;
