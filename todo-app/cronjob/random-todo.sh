#!/bin/sh

set -e

# The todo backend endpoint (can be overridden in the container)
TODO_API_URL="${TODO_API_URL:-http://localhost:8081/todos}"

# Fetch the redirect target without following it
ARTICLE_PATH=$(curl -sI https://en.wikipedia.org/wiki/Special:Random \
  | grep -i "Location:" \
  | awk '{print $2}' \
  | tr -d '\r')

# Final article URL
ARTICLE_URL="https://en.wikipedia.org${ARTICLE_PATH}"

# Build the todo JSON payload
JSON_PAYLOAD="{\"task\":\"Read ${ARTICLE_URL}\"}"

# Send to backend
curl -s -X POST \
  -H "Content-Type: application/json" \
  -d "${JSON_PAYLOAD}" \
  "${TODO_API_URL}"

echo "Created todo: Read ${ARTICLE_URL}"
