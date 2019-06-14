import React from "react";
import { formatSpan } from "./span";

export class RuleTreeView extends React.Component {
  render() {
    return (
      <div style={{ fontFamily: "monospace" }}>
        <ul>
          <li>
            <RuleTreeNode node={this.props.node} />
          </li>
        </ul>
      </div>
    )
  }
}

class RuleTreeNode extends React.Component {
  render() {
    return (
      <>
        {this.props.node.Name}
        {formatSpan(this.props.node.Span)}
        <ul>
          {(this.props.node.Children || []).map((child, idx) => (
            <li key={idx}>
              <RuleTreeNode node={child} />
            </li>
          ))}
        </ul>
      </>
    );
  }
}
