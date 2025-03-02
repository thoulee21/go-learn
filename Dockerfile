FROM golang:1.24-alpine AS builder
WORKDIR /srv/go-app
COPY . .
RUN go build -o aichatbot


FROM alpine
WORKDIR /srv/go-app
COPY --from=builder /srv/go-app/aichatbot .

CMD ["./aichatbot"]