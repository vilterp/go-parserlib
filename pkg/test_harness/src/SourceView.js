import React from "react";
import { formatSpan } from './span';
import { cursorIsWithin } from './trace';
import classNames from "classnames";
import "./GrammarView.css";
import "./SourceView.css";

// Render a syntax-highlighted view of the source from the trace.
// Highlight hovered spans.

export class SourceView extends React.Component {
  render() {
    return (
      <div className="source-view">
        {/*<SourceViewNode {...this.props} />*/}
        (source view removed)
      </div>
    );
  }
}

class SourceViewNode extends React.Component {
  render() {
    const {
      trace,
      grammar,
      highlightedSpan,
      onHighlightSpan,
      highlightedRuleID,
      onHighlightRule,
    } = this.props;

    const formattedSpan = trace ? formatSpan(trace) : null;
    const isHighlightedSpan = formattedSpan === highlightedSpan;
    const cursorWithin = trace ? cursorIsWithin(trace) : false;

    function highlightWrapper(element) {
      return (
        <span
          className={classNames("source-span", {
            highlighted: isHighlightedSpan,
          })}
          onMouseOver={() => {
            onHighlightSpan(formattedSpan, true);
            onHighlightRule(trace.RuleID, true);
          }}
          onMouseOut={() => {
            onHighlightSpan(formattedSpan, false);
            onHighlightRule(trace.RuleID, false);
          }}
        >
          {element}
        </span>
      )
    }

    if (!trace) {
      return ""; // un-filled-in sequence items
    }

    const highlightProps = {
      onHighlightSpan: onHighlightSpan,
      highlightedSpan: highlightedSpan,
      onHighlightRule: onHighlightRule,
      highlightedRuleID: highlightedRuleID,
    };

    const rule = grammar.RulesByID[trace.RuleID];
    switch (rule.RuleType) {
      case "SEQUENCE":
        return (
          <span>
            {trace.ItemTraces.map((itemTrace, idx) => (
              <SourceViewNode
                key={idx}
                trace={itemTrace}
                grammar={grammar}
                {...highlightProps}
              />
            ))}
          </span>
        );
      case "CHOICE":
        return (
          <SourceViewNode
            trace={trace.ChoiceTrace}
            grammar={grammar}
            {...highlightProps}
          />
        );
      case "REF": {
        return (
          <SourceViewNode
            trace={trace.RefTrace}
            grammar={grammar}
            {...highlightProps}
          />
        );
      }
      case "KEYWORD":
        return highlightWrapper(
          <span
            className="rule-keyword"
            style={{ textDecoration: cursorWithin ? "underline" : "none" }}
          >
            {textWithCursor(rule.Keyword, trace.CursorPos)}
          </span>
        );
      case "REGEX":
        return highlightWrapper(
          <span
            className="rule-regex"
            style={{
              whiteSpace: "pre",
              textDecoration: cursorWithin ? "underline" : "none",
            }}
          >
            {textWithCursor(trace.RegexMatch, trace.CursorPos)}
          </span>
        );
      case "MAP":
        return (
          <SourceViewNode
            trace={trace.InnerTrace}
            grammar={grammar}
            {...highlightProps}
          />
        );
      case "SUCCEED":
        return null;
      default:
        console.error(trace);
        return <pre>{JSON.stringify(trace)}</pre>
    }
  }
}

function textWithCursor(text, pos) {
  if (pos < 0 || pos >= text.length) {
    return text;
  }

  return (
    <span>
      {text.substr(0, pos)}
      <div className="cursor" />
      {text.substr(pos)}
    </span>
  );
}
