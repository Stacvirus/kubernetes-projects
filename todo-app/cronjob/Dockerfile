FROM alpine:3.19

RUN apk add --no-cache curl

WORKDIR /app
COPY random-todo.sh /app/random-todo.sh
RUN chmod +x /app/random-todo.sh

CMD ["/app/random-todo.sh"]
