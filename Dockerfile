FROM golang:1.24 AS builder
WORKDIR /srv/go-app
COPY . .
RUN go env -w GO111MODULE=on
RUN go env -w GOPROXY=https://goproxy.cn,direct
RUN go build -o microservice


FROM golang:1.24
WORKDIR /srv/go-app
COPY --from=builder /srv/go-app/microservice .

CMD ["./microservice"]