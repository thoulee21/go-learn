FROM golang:1.24-alpine AS builder
WORKDIR /srv/go-app
COPY . .

RUN apt-get update && apt-get install -y gcc libc-dev build-essential
RUN go env -w CGO_ENABLED=1

RUN go build -o aichatbot


FROM alpine
WORKDIR /srv/go-app
COPY --from=builder /srv/go-app/aichatbot .

CMD ["./aichatbot"]