FROM golang:1.24-alpine AS builder
WORKDIR /srv/go-app
COPY . .
RUN go env -w GO111MODULE=on
RUN go env -w GOPROXY=https://goproxy.cn,direct
RUN go build -o aichatbot


FROM ubuntu
WORKDIR /srv/go-app
COPY --from=builder /srv/go-app/aichatbot .

CMD ["./aichatbot"]