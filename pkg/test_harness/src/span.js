export function formatPos(pos) {
  return `${pos.Line}:${pos.Col}`
}

export function formatSpan(span) {
  return `[${formatPos(span.From)} - ${formatPos(span.To)}]`
}