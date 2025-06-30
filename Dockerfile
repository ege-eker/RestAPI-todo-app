FROM golang:1.23.8-alpine as builder

WORKDIR /app

COPY . .

RUN go build -o RestAPI-todo-app

FROM scratch as runner

WORKDIR /app

COPY --from=builder /app/RestAPI-todo-app /app/RestAPI-todo-app

CMD ["/app/RestAPI-todo-app"]