import React from 'react';

export class CompletionsView extends React.Component {
  render() {
    return (
      <ul>
        {this.props.completions.map((completion) => (
          <li key={completion.Content}>
            {completion.Kind}:
            {" "}
            <span style={{ fontFamily: "monospace" }}>
              {completion.Content}
            </span>
          </li>
        ))}
      </ul>
    )
  }
}